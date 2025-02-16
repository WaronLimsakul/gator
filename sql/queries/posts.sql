-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ( $1, NOW(), NOW(), $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostsForUser :many
SELECT * FROM posts
ORDER BY updated_at ASC
LIMIT $1;

-- name: ResetPosts :exec
DELETE FROM posts;
