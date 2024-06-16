package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"

	"github.com/gwolves/feedy/internal/channeltalk"
	"github.com/gwolves/feedy/internal/feed"
)

func newChannelTalkNotifier(
	appName string,
	client *channeltalk.Client,
	logger *slog.Logger,
) *ChannelTalkNotifier {
	return &ChannelTalkNotifier{
		appName: appName,
		client:  client,
		logger:  logger,
	}
}

type ChannelTalkNotifier struct {
	appName string
	client  *channeltalk.Client
	logger  *slog.Logger
}

func (n *ChannelTalkNotifier) Notify(
	ctx context.Context,
	channelID string,
	groupID string,
	botName string,
	blocks []channeltalk.MessageBlock,
	buttons []channeltalk.Button,
) error {
	req := channeltalk.WriteGroupMessageRequest{
		ChannelID: channelID,
		GroupID:   groupID,
		DTO: channeltalk.GroupMessage{
			BotName: botName,
			Blocks:  blocks,
			Buttons: buttons,
		},
	}

	_, err := n.client.WriteGroupMessage(ctx, &req)
	if err != nil {
		return errors.Wrap(err, "failed to send group message")
	}

	return nil
}

func (n *ChannelTalkNotifier) NotifyFeeds(
	ctx context.Context,
	channelID string,
	groupID string,
	feeds []feed.Feed,
) error {
	var blocks []channeltalk.MessageBlock
	if len(feeds) > 0 {
		bullets := make([]channeltalk.MessageBlock, 0, len(feeds))
		for _, f := range feeds {
			bullets = append(bullets, channeltalk.NewTextBlock(
				fmt.Sprintf("ID: %d - %s (%s)", f.ID, f.Name, f.URL),
			))
		}
		blocks = []channeltalk.MessageBlock{
			channeltalk.NewTextBlock("Subscriptions"),
			channeltalk.NewBulletsBlock(bullets),
		}
	} else {
		blocks = []channeltalk.MessageBlock{
			channeltalk.NewTextBlock("No Subscriptions"),
		}
	}

	return n.Notify(ctx, channelID, groupID, n.appName, blocks, nil)
}

func (n *ChannelTalkNotifier) NotifyItem(
	ctx context.Context,
	channelID string,
	groupID string,
	botName string,
	item *feed.Item,
) error {
	if botName == "" {
		botName = n.appName
	}

	blocks := []channeltalk.MessageBlock{
		channeltalk.NewTextBlock(
			channeltalk.InlineLink(item.Link, item.Title),
		),
		channeltalk.NewTextBlock(item.Content),
	}

	var buttons []channeltalk.Button
	for _, link := range item.ExtraLinks {
		buttons = append(buttons, channeltalk.Button{
			Title:        link.Value,
			ColorVariant: 1, // cobalt
			Action: struct {
				WebAction struct {
					Attributes struct {
						URL string "json:\"url\""
					} "json:\"attributes\""
				} "json:\"web_action\""
			}{
				WebAction: struct {
					Attributes struct {
						URL string "json:\"url\""
					} "json:\"attributes\""
				}{
					Attributes: struct {
						URL string "json:\"url\""
					}{
						URL: link.URL,
					},
				},
			},
		})
	}

	return n.Notify(ctx, channelID, groupID, botName, blocks, buttons)
}

func (n *ChannelTalkNotifier) NotifyString(
	ctx context.Context,
	channelID, groupID, msg string,
) error {
	blocks := []channeltalk.MessageBlock{
		channeltalk.NewTextBlock(msg),
	}

	return n.Notify(ctx, channelID, groupID, n.appName, blocks, nil)
}
