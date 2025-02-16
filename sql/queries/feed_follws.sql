-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id) 
    VALUES (
        $1,
        NOW(),
        NOW(),
        $2,
        $3
        )
    RETURNING *
) 
SELECT 
    inserted_feed_follow.*, 
    feeds.name AS feed_name, 
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowForUser :many
SELECT feed_follows.*, feeds.name AS feed_name, users.name AS user_name 
FROM feed_follows 
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows 
WHERE user_id = $1 AND feed_id = $2;

-- name: ResetFeedFollows :exec
DELETE FROM feed_follows;
