-- name: CreateCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE category_id = $1 LIMIT 1;


-- name: DeleteCategory :exec
DELETE FROM categories
WHERE category_id = $1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY category_id
LIMIT $1
OFFSET $2;

