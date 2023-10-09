-- name: CreateComment :one
INSERT INTO comments (
  comment_text, ticket_id, user_commented
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: GetCommentForUpdate :one
SELECT * FROM comments
WHERE comment_id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: ListComments :many
SELECT * FROM comments
WHERE ticket_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;


-- name: UpdateComment :one
UPDATE comments
SET comment_text = $2
WHERE comment_id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE comment_id = $1;
