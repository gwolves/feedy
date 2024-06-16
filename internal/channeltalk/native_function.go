package channeltalk

import (
	"context"
	"encoding/json"
)

func (c *Client) WriteGroupMessage(ctx context.Context, params *WriteGroupMessageRequest) (*WriteGroupMessageResponse, error) {
	token, err := c.getAccessToken(ctx, params.ChannelID)
	if err != nil {
		return nil, err
	}

	funcRes, err := c.invokeNativeFunction(ctx, token, params)
	if err != nil {
		return nil, err
	}

	var res WriteGroupMessageResponse
	if err := json.Unmarshal(funcRes, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
