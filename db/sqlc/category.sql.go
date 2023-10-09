// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: category.sql

package db

import (
	"context"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
)
RETURNING category_id, name
`

func (q *Queries) CreateCategory(ctx context.Context, name string) (Category, error) {
	row := q.db.QueryRowContext(ctx, createCategory, name)
	var i Category
	err := row.Scan(&i.CategoryID, &i.Name)
	return i, err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories
WHERE category_id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, categoryID int64) error {
	_, err := q.db.ExecContext(ctx, deleteCategory, categoryID)
	return err
}

const getCategory = `-- name: GetCategory :one
SELECT category_id, name FROM categories
WHERE category_id = $1 LIMIT 1
`

func (q *Queries) GetCategory(ctx context.Context, categoryID int64) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategory, categoryID)
	var i Category
	err := row.Scan(&i.CategoryID, &i.Name)
	return i, err
}

const listCategories = `-- name: ListCategories :many
SELECT category_id, name FROM categories
ORDER BY category_id
LIMIT $1
OFFSET $2
`

type ListCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, listCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.CategoryID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
