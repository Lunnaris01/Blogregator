-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *)
SELECT
    i.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow i
INNER JOIN users u ON u.id = i.user_id
INNER JOIN feeds f ON f.id = i.feed_id
;


-- name: GetFeedFollowsForUser :many
SELECT f.id,f.url,f.name as "Feedname" FROM feed_follows ff 
INNER JOIN users u ON ff.user_id = u.id 
INNER JOIN feeds f on ff.feed_id = f.id
WHERE u.id = $1;

-- name: UnfollowFeed :one
WITH deleted_follow AS(
DELETE FROM feed_follows USING feeds 
WHERE feeds.ID = feed_follows.feed_id 
AND feeds.url = $1 
AND feed_follows.user_id = $2
RETURNING feed_follows.feed_id)
SELECT feeds.name FROM feeds INNER JOIN deleted_follow
ON deleted_follow.feed_id = feeds.id;