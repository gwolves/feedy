// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuditLog = `-- name: CreateAuditLog :exec
INSERT INTO audit_logs (actor, action, target)
VALUES ($1, $2, $3)
`

type CreateAuditLogParams struct {
	Actor  string
	Action string
	Target pgtype.Text
}

func (q *Queries) CreateAuditLog(ctx context.Context, arg CreateAuditLogParams) error {
	_, err := q.db.Exec(ctx, createAuditLog, arg.Actor, arg.Action, arg.Target)
	return err
}

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (name, url) VALUES ($1, $2)
RETURNING id, name, url, created_at
`

type CreateFeedParams struct {
	Name string
	Url  string
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRow(ctx, createFeed, arg.Name, arg.Url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}

const createSubscription = `-- name: CreateSubscription :one
INSERT INTO subscriptions (bot_name, feed_id, channel_id, group_id)
VALUES ($1, $2, $3, $4)
RETURNING id, bot_name, feed_id, channel_id, group_id, published_at, created_at
`

type CreateSubscriptionParams struct {
	BotName   pgtype.Text
	FeedID    int64
	ChannelID string
	GroupID   string
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRow(ctx, createSubscription,
		arg.BotName,
		arg.FeedID,
		arg.ChannelID,
		arg.GroupID,
	)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.BotName,
		&i.FeedID,
		&i.ChannelID,
		&i.GroupID,
		&i.PublishedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteSubscription = `-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE channel_id = $1
  AND group_id = $2
  AND feed_id = $3
`

type DeleteSubscriptionParams struct {
	ChannelID string
	GroupID   string
	FeedID    int64
}

func (q *Queries) DeleteSubscription(ctx context.Context, arg DeleteSubscriptionParams) error {
	_, err := q.db.Exec(ctx, deleteSubscription, arg.ChannelID, arg.GroupID, arg.FeedID)
	return err
}

const getFeedByID = `-- name: GetFeedByID :one
SELECT id, name, url, created_at FROM feeds
WHERE id = $1
`

func (q *Queries) GetFeedByID(ctx context.Context, id int64) (Feed, error) {
	row := q.db.QueryRow(ctx, getFeedByID, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}

const getFeedByURL = `-- name: GetFeedByURL :one
SELECT id, name, url, created_at FROM feeds
WHERE url = $1
`

func (q *Queries) GetFeedByURL(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRow(ctx, getFeedByURL, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}

const getSubscription = `-- name: GetSubscription :one
SELECT id, bot_name, feed_id, channel_id, group_id, published_at, created_at FROM subscriptions
WHERE feed_id = $1
  AND channel_id = $2
  AND group_id = $3
`

type GetSubscriptionParams struct {
	FeedID    int64
	ChannelID string
	GroupID   string
}

func (q *Queries) GetSubscription(ctx context.Context, arg GetSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRow(ctx, getSubscription, arg.FeedID, arg.ChannelID, arg.GroupID)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.BotName,
		&i.FeedID,
		&i.ChannelID,
		&i.GroupID,
		&i.PublishedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listFeeds = `-- name: ListFeeds :many
SELECT id, name, url, created_at FROM feeds
ORDER BY id
`

func (q *Queries) ListFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.Query(ctx, listFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSubscribedFeedsByGroup = `-- name: ListSubscribedFeedsByGroup :many
SELECT
  f.id, f.name, f.url, f.created_at
FROM subscriptions s
  INNER JOIN feeds f on s.feed_id = f.id
WHERE
  s.channel_id = $1
  AND s.group_id = $2
ORDER BY f.id
`

type ListSubscribedFeedsByGroupParams struct {
	ChannelID string
	GroupID   string
}

func (q *Queries) ListSubscribedFeedsByGroup(ctx context.Context, arg ListSubscribedFeedsByGroupParams) ([]Feed, error) {
	rows, err := q.db.Query(ctx, listSubscribedFeedsByGroup, arg.ChannelID, arg.GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSubscriptionsByFeed = `-- name: ListSubscriptionsByFeed :many
SELECT id, bot_name, feed_id, channel_id, group_id, published_at, created_at FROM subscriptions
WHERE feed_id = $1
ORDER BY id
`

func (q *Queries) ListSubscriptionsByFeed(ctx context.Context, feedID int64) ([]Subscription, error) {
	rows, err := q.db.Query(ctx, listSubscriptionsByFeed, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.BotName,
			&i.FeedID,
			&i.ChannelID,
			&i.GroupID,
			&i.PublishedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSubscriptionPublishedAt = `-- name: UpdateSubscriptionPublishedAt :exec
UPDATE subscriptions SET published_at = $1
WHERE id = $2
`

type UpdateSubscriptionPublishedAtParams struct {
	PublishedAt pgtype.Timestamptz
	ID          int64
}

func (q *Queries) UpdateSubscriptionPublishedAt(ctx context.Context, arg UpdateSubscriptionPublishedAtParams) error {
	_, err := q.db.Exec(ctx, updateSubscriptionPublishedAt, arg.PublishedAt, arg.ID)
	return err
}
