-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ResetFeeds :exec
DELETE FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name AS "username" 
FROM feeds f INNER JOIN users u ON f.user_id = u.id;
