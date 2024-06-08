-- name: GetFeedByID :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY id;

-- name: CreateFeed :one
INSERT INTO feeds (name, url) VALUES ($1, $2)
RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE feed_id = $1
  AND channel_id = $2
  AND group_id = $3;

-- name: ListSubscriptionsByFeed :many
SELECT * FROM subscriptions
WHERE feed_id = $1
ORDER BY id;

-- name: ListSubscribedFeedsByGroup :many
SELECT
  f.*
FROM subscriptions s
  INNER JOIN feeds f on s.feed_id = f.id
WHERE
  s.channel_id = $1
  AND s.group_id = $2
ORDER BY f.id;

-- name: CreateSubscription :one
INSERT INTO subscriptions (bot_name, feed_id, channel_id, group_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateSubscriptionPublishedAt :exec
UPDATE subscriptions SET published_at = $1
WHERE id = $2;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE channel_id = $1
  AND group_id = $2
  AND feed_id = $3;

-- name: CreateAuditLog :exec
INSERT INTO audit_logs (actor, action, target)
VALUES ($1, $2, $3);
