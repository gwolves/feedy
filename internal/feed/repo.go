package feed

import (
	"context"
	"time"
)

type Repository interface {
	WithUnitOfWork(context.Context) (UnitOfWork, Repository, error)

	GetFeedByID(ctx context.Context, id int64) (*Feed, error)
	GetFeedByURL(ctx context.Context, url string) (*Feed, error)
	ListFeeds(ctx context.Context) ([]Feed, error)
	CreateFeed(context.Context, *Feed) (*Feed, error)

	ListSubscriptionsByFeed(ctx context.Context, feedID int64) ([]Subscription, error)
	ListSubscribedFeedsByGroup(
		ctx context.Context,
		channelID string,
		groupID string,
	) ([]Feed, error)
	CreateSubscription(context.Context, *Subscription) (*Subscription, error)
	DeleteSubscription(
		ctx context.Context,
		channelID string,
		groupID string,
		feedID int64,
	) error
	TouchSubscription(context.Context, *Subscription, time.Time) error

	CreateAuditLog(cxt context.Context, actor, action, target string) error
}

type UnitOfWork interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
