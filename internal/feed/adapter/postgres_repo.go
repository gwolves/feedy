package adapter

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gwolves/feedy/internal/feed"
	"github.com/gwolves/feedy/internal/sql"
)

func NewPostgresRepo(conn *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{
		conn:    conn,
		queries: sql.New(conn),
	}
}

type PostgresRepo struct {
	conn    *pgx.Conn
	queries *sql.Queries
}

func (r *PostgresRepo) WithUnitOfWork(ctx context.Context) (feed.UnitOfWork, feed.Repository, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	uow := &transaction{
		tx: tx,
	}

	repo := &PostgresRepo{
		conn:    tx.Conn(),
		queries: r.queries.WithTx(tx),
	}

	return uow, repo, nil
}

func (r *PostgresRepo) GetFeedByID(ctx context.Context, id int64) (*feed.Feed, error) {
	dto, err := r.queries.GetFeedByID(ctx, id)
	if err != nil {
		// TODO:
		// return nil, err
		return nil, nil
	}

	return &feed.Feed{
		ID:   dto.ID,
		Name: dto.Name,
		URL:  dto.Url,
	}, nil
}

func (r *PostgresRepo) GetFeedByURL(ctx context.Context, url string) (*feed.Feed, error) {
	dto, err := r.queries.GetFeedByURL(ctx, url)
	if err != nil {
		// TODO:
		// return nil, err
		return nil, nil
	}

	return &feed.Feed{
		ID:   dto.ID,
		Name: dto.Name,
		URL:  dto.Url,
	}, nil
}

func (r *PostgresRepo) ListFeeds(ctx context.Context) ([]feed.Feed, error) {
	dtos, err := r.queries.ListFeeds(ctx)
	if err != nil {
		return nil, err
	}

	var feeds []feed.Feed
	if len(dtos) > 0 {
		feeds = make([]feed.Feed, 0, len(dtos))
		for _, dto := range dtos {
			feeds = append(feeds, feed.Feed{
				ID:   dto.ID,
				Name: dto.Name,
				URL:  dto.Url,
			})
		}
	}

	return feeds, nil
}

func (r *PostgresRepo) CreateFeed(ctx context.Context, f *feed.Feed) (*feed.Feed, error) {
	dto, err := r.queries.CreateFeed(ctx, sql.CreateFeedParams{
		Name: f.Name,
		Url:  f.URL,
	})
	if err != nil {
		return nil, err
	}

	return &feed.Feed{
		ID:   dto.ID,
		Name: dto.Name,
		URL:  dto.Url,
	}, nil
}

func (r *PostgresRepo) ListSubscriptionsByFeed(ctx context.Context, feedID int64) ([]feed.Subscription, error) {
	dtos, err := r.queries.ListSubscriptionsByFeed(ctx, feedID)
	if err != nil {
		return nil, err
	}

	var subs []feed.Subscription
	if len(dtos) > 0 {
		subs = make([]feed.Subscription, 0, len(dtos))
		for _, dto := range dtos {
			subs = append(subs, feed.Subscription{
				ID:          dto.ID,
				ChannelID:   dto.ChannelID,
				GroupID:     dto.GroupID,
				FeedID:      dto.FeedID,
				BotName:     dto.BotName.String,
				PublishedAt: dto.PublishedAt.Time,
			})
		}
	}

	return subs, nil
}

func (r *PostgresRepo) ListSubscribedFeedsByGroup(
	ctx context.Context,
	channelID string,
	groupID string,
) ([]feed.Feed, error) {
	dtos, err := r.queries.ListSubscribedFeedsByGroup(
		ctx,
		sql.ListSubscribedFeedsByGroupParams{
			ChannelID: channelID,
			GroupID:   groupID,
		},
	)
	if err != nil {
		return nil, err
	}

	var feeds []feed.Feed
	if len(dtos) > 0 {
		feeds = make([]feed.Feed, 0, len(dtos))
		for _, dto := range dtos {
			feeds = append(feeds, feed.Feed{
				ID:   dto.ID,
				Name: dto.Name,
				URL:  dto.Url,
			})
		}
	}

	return feeds, nil
}

func (r *PostgresRepo) CreateSubscription(
	ctx context.Context,
	sub *feed.Subscription,
) (*feed.Subscription, error) {
	dto, err := r.queries.CreateSubscription(ctx, sql.CreateSubscriptionParams{
		BotName: pgtype.Text{
			String: sub.BotName,
			Valid:  sub.BotName != "",
		},
		FeedID:    sub.FeedID,
		ChannelID: sub.ChannelID,
		GroupID:   sub.GroupID,
	})
	if err != nil {
		return nil, err
	}

	return &feed.Subscription{
		ID:          dto.ID,
		ChannelID:   dto.ChannelID,
		GroupID:     dto.GroupID,
		FeedID:      dto.FeedID,
		BotName:     dto.BotName.String,
		PublishedAt: dto.PublishedAt.Time,
	}, nil
}

func (r *PostgresRepo) DeleteSubscription(
	ctx context.Context,
	channelID string,
	groupID string,
	feedID int64,
) error {
	return r.queries.DeleteSubscription(ctx, sql.DeleteSubscriptionParams{
		ChannelID: channelID,
		GroupID:   groupID,
		FeedID:    feedID,
	})
}

func (r *PostgresRepo) TouchSubscription(
	ctx context.Context,
	sub *feed.Subscription,
	time time.Time,
) error {
	return r.queries.UpdateSubscriptionPublishedAt(ctx, sql.UpdateSubscriptionPublishedAtParams{
		ID: sub.ID,
		PublishedAt: pgtype.Timestamptz{
			Time:  time,
			Valid: true,
		},
	})
}

func (r *PostgresRepo) CreateAuditLog(
	ctx context.Context,
	actor string,
	action string,
	target string,
) error {
	return r.queries.CreateAuditLog(ctx, sql.CreateAuditLogParams{
		Actor:  actor,
		Action: action,
		Target: pgtype.Text{
			String: target,
			Valid:  target != "",
		},
	})
}

type transaction struct {
	tx pgx.Tx
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
