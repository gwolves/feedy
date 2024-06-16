package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/gwolves/feedy/internal/channeltalk"
	"github.com/gwolves/feedy/internal/feed"
)

func NewUseCase(
	appName string,
	repo feed.Repository,
	client *channeltalk.Client,
	logger *slog.Logger,
) *UseCase {
	notifier := newChannelTalkNotifier(appName, client, logger)
	return &UseCase{
		parser:   gofeed.NewParser(),
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

type UseCase struct {
	appName  string
	parser   *gofeed.Parser
	repo     feed.Repository
	notifier *ChannelTalkNotifier
	logger   *slog.Logger
}

func (u *UseCase) ListSubscribedFeeds(
	ctx context.Context,
	channelID string,
	groupID string,
	notify bool,
) ([]feed.Feed, error) {
	feeds, err := u.repo.ListSubscribedFeedsByGroup(ctx, channelID, groupID)
	if err != nil {
		return nil, err
	}

	if !notify {
		return feeds, nil
	}

	if err = u.notifier.NotifyFeeds(ctx, channelID, groupID, feeds); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (u *UseCase) Subscribe(
	ctx context.Context,
	channelID string,
	groupID string,
	url string,
	botName string,
) (err error) {
	uow, repo, err := u.repo.WithUnitOfWork(ctx)
	defer uow.Rollback(ctx)

	if err != nil {
		return err
	}

	f, err := repo.GetFeedByURL(ctx, url)
	if err != nil {
		return errors.Wrap(err, "failed to get feed")
	}

	if f == nil {
		f, err = u.createFeed(ctx, url, repo)
		if err != nil {
			return WithReason(err, fmt.Sprintf("Invalid Feed: %s", url))
		}
	}

	sub, err := repo.CreateSubscription(ctx, &feed.Subscription{
		FeedID:    f.ID,
		ChannelID: channelID,
		GroupID:   groupID,
		BotName:   botName,
	})
	if err != nil {
		return WithReason(err, "Failed to subscribe")
	}

	if err = repo.CreateAuditLog(ctx, getCaller(ctx), "subscribe", fmt.Sprintf("sub:%d", sub.ID)); err != nil {
		return err
	}

	if err = uow.Commit(ctx); err != nil {
		return err
	}

	return u.notifier.NotifyString(
		ctx,
		channelID,
		groupID,
		fmt.Sprintf("Subscribed: %s (%s)", f.Name, f.URL),
	)
}

func (u *UseCase) createFeed(ctx context.Context, url string, repo feed.Repository) (*feed.Feed, error) {
	f, err := u.parser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	feed := feed.Feed{
		Name: f.Title,
		URL:  url,
	}
	return repo.CreateFeed(ctx, &feed)
}

func (u *UseCase) Unsubscribe(
	ctx context.Context,
	channelID string,
	groupID string,
	feedID int64,
) error {
	f, err := u.repo.GetFeedByID(ctx, feedID)
	if err != nil {
		return WithReason(err, fmt.Sprintf("Failed to get feed %d", feedID))
	}

	if f == nil {
		return WithReason(err, fmt.Sprintf("No subscription for feed: %d", feedID))
	}

	uow, repo, err := u.repo.WithUnitOfWork(ctx)
	defer uow.Rollback(ctx)

	if err := repo.DeleteSubscription(ctx, channelID, groupID, feedID); err != nil {
		return WithReason(err, fmt.Sprintf("No subscription for feed: %d", feedID))
	}

	if err = repo.CreateAuditLog(
		ctx,
		getCaller(ctx),
		"unsubscribe",
		fmt.Sprintf("sub:%s:%s:%d", channelID, groupID, feedID),
	); err != nil {
		return err
	}

	if err := uow.Commit(ctx); err != nil {
		return err
	}

	return u.notifier.NotifyString(ctx, channelID, groupID, fmt.Sprintf("Unsubscribed: %s (%s)", f.Name, f.URL))
}

func (u *UseCase) PublishFeed(ctx context.Context, feedID int64) error {
	f, err := u.repo.GetFeedByID(ctx, feedID)
	if err != nil {
		return err
	}

	if f == nil {
		return errors.Errorf("feed not exist: %d", feedID)
	}

	return u.publishFeed(ctx, f)
}

func (u *UseCase) PublishAllFeeds(ctx context.Context) error {
	feeds, err := u.repo.ListFeeds(ctx)
	if err != nil {
		return err
	}

	for _, f := range feeds {
		err = u.publishFeed(ctx, &f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UseCase) publishFeed(ctx context.Context, f *feed.Feed) error {
	subs, err := u.repo.ListSubscriptionsByFeed(ctx, f.ID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		u.logger.Debug("no subscription")
		return nil
	}

	u.logger.Info("fetch start", "feed_id", f.ID, "feed_name", f.Name)
	fetcher := feed.NewFetcher(u.logger)
	items, err := fetcher.Fetch(f)
	if err != nil {
		return err
	}
	u.logger.Info("fetch end", "count", len(items))

	for _, sub := range subs {
		u.logger.Info("publish start", "subscription_id", sub.ID, "last_published_at", sub.PublishedAt)

		var lastPublished *time.Time
		for _, item := range items {
			if !sub.PublishedAt.Before(item.PublishedAt) {
				u.logger.Debug("already published item", "title", item.Title)
				continue
			}
			u.logger.Debug("item", "title", item.Title)

			err = u.notifier.NotifyItem(ctx, sub.ChannelID, sub.GroupID, sub.BotName, &item)
			if err != nil {
				u.logger.Error("notification failed", "error", err)
				continue
			}

			if lastPublished == nil || lastPublished.Before(item.PublishedAt) {
				lastPublished = &item.PublishedAt
			}
		}

		if lastPublished != nil {
			err = u.repo.TouchSubscription(ctx, &sub, *lastPublished)
			if err != nil {
				u.logger.Error("touch failed", "error", err)
				continue
			}
			u.logger.Info("publish done", "published_at", *lastPublished)
		}
	}

	// TODO: report error
	return nil
}

func (u *UseCase) Notify(
	ctx context.Context,
	channelID, groupID, msg string,
) error {
	return u.notifier.NotifyString(ctx, channelID, groupID, msg)
}

func getCaller(ctx context.Context) string {
	if v := ctx.Value("caller"); v != nil {
		return v.(string)
	}

	return "unknown"
}
