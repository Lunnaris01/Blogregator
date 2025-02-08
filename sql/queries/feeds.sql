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


-- name: MarkFeedFetched :exec
UPDATE feeds SET updated_at = $2, last_fetched_at = $3  WHERE id =$1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;