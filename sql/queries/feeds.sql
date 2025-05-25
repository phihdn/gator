-- name: CreateFeed :one
INSERT INTO
    feeds (id, created_at, updated_at, name, url, user_id)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetFeedsByUser :many
SELECT
    *
FROM
    feeds
WHERE
    user_id = $1
ORDER BY
    created_at DESC;

-- name: GetAllFeedsWithUsers :many
SELECT
    f.id, f.created_at, f.updated_at, f.name, f.url, f.user_id, u.name AS user_name
FROM
    feeds f
JOIN
    users u ON f.user_id = u.id
ORDER BY
    f.created_at DESC;
