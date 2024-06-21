package channeltalk

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

const (
	defaultEndpoint    = "https://app-store-api.channel.io"
	nativeFunctionPath = "/general/v1/native/functions"

	refreshTokenExpiry = 7 * 24 * time.Hour
)

func NewClient(secret string, logger *slog.Logger, opts ...Option) *Client {
	var config config
	for _, opt := range opts {
		opt(&config)
	}

	if config.endpoint == "" {
		config.endpoint = defaultEndpoint
	}

	if config.httpClient == nil {
		config.httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Client{
		secret: secret,
		client: resty.
			NewWithClient(config.httpClient).
			SetBaseURL(config.endpoint),
		cache:  cache.New(refreshTokenExpiry, 1*time.Hour),
		logger: logger,
	}
}

type Client struct {
	client *resty.Client
	secret string
	cache  *cache.Cache
	logger *slog.Logger
}

type Option func(*config)

type config struct {
	endpoint   string
	httpClient *http.Client
}

func WithEndpoint(endpoint string) Option {
	return func(o *config) {
		o.endpoint = endpoint
	}
}

func (c *Client) invokeNativeFunction(
	ctx context.Context,
	accessToken string,
	params nativeFuntcionParams,
) (json.RawMessage, error) {
	body := nativeFunctionRequest{
		Method: params.Method(),
		Params: params,
	}

	r := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetJSONEscapeHTML(false).
		SetBody(&body).
		SetResult(&nativeFunctionResponse{}).
		SetError(&nativeFunctionErrorResponse{})

	if accessToken != "" {
		r = r.SetHeader("x-access-token", accessToken)
	}

	resp, err := r.Put(nativeFunctionPath)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.Errorf("error code: %d. %s", resp.StatusCode(), resp.Error())
	}

	res := resp.Result().(*nativeFunctionResponse)
	if res.Error != nil {
		return nil, res.Error
	}

	return res.Result, nil
}

func (c *Client) issueToken(ctx context.Context, channelID *string) (*IssueTokenResponse, error) {
	funcRes, err := c.invokeNativeFunction(ctx, "", &issueTokenRequest{
		Secret:    c.secret,
		ChannelId: channelID,
	})
	if err != nil {
		return nil, err
	}

	var res IssueTokenResponse
	if err = json.Unmarshal(funcRes, &res); err != nil {
		return nil, err
	}

	res.Expiry = time.Now().Add(time.Duration(res.ExpiresIn)*time.Second - 20*time.Second) // with buffer
	return &res, nil
}

func (c *Client) refreshIssueToken(ctx context.Context, refreshToken string) (*IssueTokenResponse, error) {
	funcRes, err := c.invokeNativeFunction(ctx, "", &refreshIssueTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	var res IssueTokenResponse
	if err = json.Unmarshal(funcRes, &res); err != nil {
		return nil, err
	}

	res.Expiry = time.Now().Add(time.Duration(res.ExpiresIn)*time.Second - 20*time.Second) // with buffer
	return &res, nil
}

func (c *Client) getAccessToken(ctx context.Context, channelID string) (string, error) {
	// 1. from local cache
	cacheKey := fmt.Sprintf("auth:%s", channelID)
	token, ok := c.cache.Get(cacheKey)

	if ok {
		token := token.(IssueTokenResponse)
		if time.Now().Before(token.Expiry) {
			c.logger.Debug("access_token", "from", "cached")
			return token.AccessToken, nil
		}

		// 2. refresh token
		res, err := c.refreshIssueToken(ctx, token.RefreshToken)
		if err == nil {
			c.logger.Debug("access_token", "from", "refresh_token")
			c.cache.Set(cacheKey, *res, refreshTokenExpiry)
			return res.AccessToken, nil
		}
	}

	// 3. issue token
	res, err := c.issueToken(ctx, &channelID)
	if err != nil {
		return "", err
	}
	c.logger.Debug("access_token", "from", "issue_token")

	c.cache.Set(cacheKey, *res, refreshTokenExpiry)
	return res.AccessToken, nil
}
