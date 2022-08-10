// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: transaction_type.sql

package db

import (
	"context"
)

const allTransactionType = `-- name: AllTransactionType :many
SELECT id, type_name, type_description, created_at, updated_at, deleted_at
FROM od_transaction_types
ORDER BY created_at
LIMIT $1
OFFSET $2
`

type AllTransactionTypeParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) AllTransactionType(ctx context.Context, arg AllTransactionTypeParams) ([]OdTransactionType, error) {
	rows, err := q.db.QueryContext(ctx, allTransactionType, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OdTransactionType
	for rows.Next() {
		var i OdTransactionType
		if err := rows.Scan(
			&i.ID,
			&i.TypeName,
			&i.TypeDescription,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
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

const createTransactionType = `-- name: CreateTransactionType :one
INSERT INTO od_transaction_types (
    type_name, type_description
) VALUES (
    $1, $2
)
RETURNING id, type_name, type_description, created_at, updated_at, deleted_at
`

type CreateTransactionTypeParams struct {
	TypeName        string `json:"type_name"`
	TypeDescription string `json:"type_description"`
}

func (q *Queries) CreateTransactionType(ctx context.Context, arg CreateTransactionTypeParams) (OdTransactionType, error) {
	row := q.db.QueryRowContext(ctx, createTransactionType, arg.TypeName, arg.TypeDescription)
	var i OdTransactionType
	err := row.Scan(
		&i.ID,
		&i.TypeName,
		&i.TypeDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteTransactionType = `-- name: DeleteTransactionType :exec
DELETE FROM od_transaction_types 
WHERE id = $1
`

func (q *Queries) DeleteTransactionType(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransactionType, id)
	return err
}

const updateTransactionType = `-- name: UpdateTransactionType :one
UPDATE od_transaction_types 
SET 
    type_name = CASE WHEN $2::text <> '' THEN $2::text ELSE type_name END,
    type_description = CASE WHEN $3::text <> '' THEN $3::text ELSE type_description END,
    updated_at = now()
WHERE id = $1
RETURNING id, type_name, type_description, created_at, updated_at, deleted_at
`

type UpdateTransactionTypeParams struct {
	ID              int64  `json:"id"`
	TypeName        string `json:"type_name"`
	TypeDescription string `json:"type_description"`
}

func (q *Queries) UpdateTransactionType(ctx context.Context, arg UpdateTransactionTypeParams) (OdTransactionType, error) {
	row := q.db.QueryRowContext(ctx, updateTransactionType, arg.ID, arg.TypeName, arg.TypeDescription)
	var i OdTransactionType
	err := row.Scan(
		&i.ID,
		&i.TypeName,
		&i.TypeDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const viewTransactionType = `-- name: ViewTransactionType :one
SELECT id, type_name, type_description, created_at, updated_at, deleted_at
FROM od_transaction_types
WHERE id = $1 LIMIT 1
`

func (q *Queries) ViewTransactionType(ctx context.Context, id int64) (OdTransactionType, error) {
	row := q.db.QueryRowContext(ctx, viewTransactionType, id)
	var i OdTransactionType
	err := row.Scan(
		&i.ID,
		&i.TypeName,
		&i.TypeDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
