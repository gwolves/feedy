package http

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/gwolves/feedy/internal/service"
)

const (
	subscribe         = "subscribe"
	unsubscribe       = "unsubscribe"
	listSubscriptions = "listSubscriptions"

	autoCompleteUnsubscribe = "autoCompleteUnsubscribe"
)

type functionHandler struct {
	u      *service.UseCase
	logger *slog.Logger
}

func (h *functionHandler) Handle(ctx context.Context, req *functionRequest) (*functionResponse, error) {
	h.logger.Debug("function request", "request", req)
	ctx = context.WithValue(
		ctx,
		"caller",
		fmt.Sprintf("%s:%s", req.Context.Caller.Type, req.Context.Caller.ID),
	)

	// currently only allow invocation from group chat
	if req.Params.Chat.Type != "group" {
		return nil, errors.Errorf("not allowed type: %s", req.Params.Chat.Type)
	}

	var res *functionResponse
	var err error

	switch req.Method {
	case subscribe:
		res, err = h.handleSubscribe(ctx, req)

	case unsubscribe:
		res, err = h.handleUnsubscribe(ctx, req)

	case listSubscriptions:
		res, err = h.handleListSubscriptions(ctx, req)

	case autoCompleteUnsubscribe:
		res, err = h.handleAutoCompleteUnubscribe(ctx, req)

	default:
		return nil, errors.New("not supported")
	}

	if err != nil {
		var expectedErr service.ExpectedError
		if errors.As(err, &expectedErr) {
			h.logger.Debug("expected error", "reason", expectedErr.Reason(), "error", err)
			h.u.Notify(
				ctx,
				req.Context.Channel.ID,
				req.Params.Chat.ID,
				expectedErr.Reason(),
			)
			return &succeedResponse, nil
		}

		return nil, err
	}

	return res, nil
}

func (h *functionHandler) handleSubscribe(ctx context.Context, req *functionRequest) (*functionResponse, error) {
	var input subscribeInputs
	if err := json.Unmarshal(req.Params.Input, &input); err != nil {
		return nil, err
	}

	channelID := req.Context.Channel.ID
	groupID := req.Params.Chat.ID

	if err := h.u.Subscribe(ctx, channelID, groupID, input.Url, input.BotName); err != nil {
		return nil, err
	}

	return &succeedResponse, nil
}

func (h *functionHandler) handleUnsubscribe(ctx context.Context, req *functionRequest) (*functionResponse, error) {
	var input unsubscribeInputs
	if err := json.Unmarshal(req.Params.Input, &input); err != nil {
		return nil, err
	}

	channelID := req.Context.Channel.ID
	groupID := req.Params.Chat.ID

	if err := h.u.Unsubscribe(ctx, channelID, groupID, input.ID); err != nil {
		return nil, err
	}

	return &succeedResponse, nil
}

func (h *functionHandler) handleListSubscriptions(ctx context.Context, req *functionRequest) (*functionResponse, error) {
	channelID := req.Context.Channel.ID
	groupID := req.Params.Chat.ID

	if _, err := h.u.ListSubscribedFeeds(ctx, channelID, groupID, true); err != nil {
		return nil, err
	}

	return &succeedResponse, nil
}

func (h *functionHandler) handleAutoCompleteUnubscribe(ctx context.Context, req *functionRequest) (*functionResponse, error) {
	channelID := req.Context.Channel.ID
	groupID := req.Params.Chat.ID

	feeds, err := h.u.ListSubscribedFeeds(ctx, channelID, groupID, false)
	if err != nil {
		return nil, err
	}

	choices := make([]choice, 0, len(feeds))
	for _, f := range feeds {
		choices = append(choices, choice{
			Name:  f.Name,
			Value: fmt.Sprintf("%d", f.ID),
		})
	}

	return &functionResponse{
		Result: autoCompleteUnsubscribeResponse{
			Choices: choices,
		},
	}, nil
}
