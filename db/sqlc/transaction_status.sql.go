// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: transaction_status.sql

package db

import (
	"context"
)

const allTransactionStatus = `-- name: AllTransactionStatus :many
SELECT id, status_name, status_description, created_at, updated_at, deleted_at
FROM transaction_statuses
ORDER BY created_at
OFFSET $1
LIMIT $2
`

type AllTransactionStatusParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) AllTransactionStatus(ctx context.Context, arg AllTransactionStatusParams) ([]TransactionStatus, error) {
	rows, err := q.db.QueryContext(ctx, allTransactionStatus, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TransactionStatus{}
	for rows.Next() {
		var i TransactionStatus
		if err := rows.Scan(
			&i.ID,
			&i.StatusName,
			&i.StatusDescription,
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

const createTransactionStatus = `-- name: CreateTransactionStatus :one
INSERT INTO transaction_statuses (
    status_name, status_description
) VALUES (
    $1, $2
)
RETURNING id, status_name, status_description, created_at, updated_at, deleted_at
`

type CreateTransactionStatusParams struct {
	StatusName        string `json:"status_name"`
	StatusDescription string `json:"status_description"`
}

func (q *Queries) CreateTransactionStatus(ctx context.Context, arg CreateTransactionStatusParams) (TransactionStatus, error) {
	row := q.db.QueryRowContext(ctx, createTransactionStatus, arg.StatusName, arg.StatusDescription)
	var i TransactionStatus
	err := row.Scan(
		&i.ID,
		&i.StatusName,
		&i.StatusDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteTransactionStatus = `-- name: DeleteTransactionStatus :exec
DELETE FROM transaction_statuses
WHERE id = $1
`

func (q *Queries) DeleteTransactionStatus(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransactionStatus, id)
	return err
}

const updateTransactionStatus = `-- name: UpdateTransactionStatus :one
UPDATE transaction_statuses 
SET 
    status_name = CASE WHEN $2::text <> '' THEN $2::text ELSE status_name END,
    status_description = CASE WHEN $3::text <> '' THEN $3::text ELSE status_description END,
    updated_at = now()
WHERE id = $1
RETURNING id, status_name, status_description, created_at, updated_at, deleted_at
`

type UpdateTransactionStatusParams struct {
	ID                int64  `json:"id"`
	StatusName        string `json:"status_name"`
	StatusDescription string `json:"status_description"`
}

func (q *Queries) UpdateTransactionStatus(ctx context.Context, arg UpdateTransactionStatusParams) (TransactionStatus, error) {
	row := q.db.QueryRowContext(ctx, updateTransactionStatus, arg.ID, arg.StatusName, arg.StatusDescription)
	var i TransactionStatus
	err := row.Scan(
		&i.ID,
		&i.StatusName,
		&i.StatusDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const viewTransactionStatus = `-- name: ViewTransactionStatus :one
SELECT id, status_name, status_description, created_at, updated_at, deleted_at
FROM transaction_statuses
WHERE id = $1 LIMIT 1
`

func (q *Queries) ViewTransactionStatus(ctx context.Context, id int64) (TransactionStatus, error) {
	row := q.db.QueryRowContext(ctx, viewTransactionStatus, id)
	var i TransactionStatus
	err := row.Scan(
		&i.ID,
		&i.StatusName,
		&i.StatusDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
