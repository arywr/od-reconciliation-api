// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: merhcant_transaction.sql

package db

import (
	"context"
	"time"

	"github.com/gobuffalo/nulls"
)

const allDuplicateMerchantTrx = `-- name: AllDuplicateMerchantTrx :many
SELECT 
    p.id, p.merchant_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM (
    SELECT id, transaction_status_id, transaction_type_id, progress_event_id, merchant_transaction_id, owner_id, transaction_id, transaction_date, transaction_datetime, collected_amount, settled_amount, created_at, updated_at, deleted_at,
    ROW_NUMBER() OVER(PARTITION BY merchant_transaction_id, transaction_date ORDER BY id asc) AS res
    FROM merchant_transactions
) p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
where p.res > 1
OFFSET $1
LIMIT $2
`

type AllDuplicateMerchantTrxParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type AllDuplicateMerchantTrxRow struct {
	ID                    int64       `json:"id"`
	MerchantTransactionID string      `json:"merchant_transaction_id"`
	OwnerID               string      `json:"owner_id"`
	TransactionID         string      `json:"transaction_id"`
	TransactionDate       time.Time   `json:"transaction_date"`
	TransactionDatetime   time.Time   `json:"transaction_datetime"`
	CollectedAmount       float64     `json:"collected_amount"`
	SettledAmount         float64     `json:"settled_amount"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	DeletedAt             nulls.Time  `json:"deleted_at"`
	ProgressEventID       nulls.Int32 `json:"progress_event_id"`
	TransactionStatusID   int16       `json:"transaction_status_id"`
	StatusName            string      `json:"status_name"`
	StatusDescription     string      `json:"status_description"`
	TransactionTypeID     int16       `json:"transaction_type_id"`
	TypeName              string      `json:"type_name"`
	TypeDescription       string      `json:"type_description"`
}

func (q *Queries) AllDuplicateMerchantTrx(ctx context.Context, arg AllDuplicateMerchantTrxParams) ([]AllDuplicateMerchantTrxRow, error) {
	rows, err := q.db.QueryContext(ctx, allDuplicateMerchantTrx, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AllDuplicateMerchantTrxRow{}
	for rows.Next() {
		var i AllDuplicateMerchantTrxRow
		if err := rows.Scan(
			&i.ID,
			&i.MerchantTransactionID,
			&i.OwnerID,
			&i.TransactionID,
			&i.TransactionDate,
			&i.TransactionDatetime,
			&i.CollectedAmount,
			&i.SettledAmount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ProgressEventID,
			&i.TransactionStatusID,
			&i.StatusName,
			&i.StatusDescription,
			&i.TransactionTypeID,
			&i.TypeName,
			&i.TypeDescription,
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

const allMerchantTrx = `-- name: AllMerchantTrx :many
SELECT 
    p.id, p.merchant_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM merchant_transactions p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
OFFSET $1
LIMIT $2
`

type AllMerchantTrxParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type AllMerchantTrxRow struct {
	ID                    int64       `json:"id"`
	MerchantTransactionID string      `json:"merchant_transaction_id"`
	OwnerID               string      `json:"owner_id"`
	TransactionID         string      `json:"transaction_id"`
	TransactionDate       time.Time   `json:"transaction_date"`
	TransactionDatetime   time.Time   `json:"transaction_datetime"`
	CollectedAmount       float64     `json:"collected_amount"`
	SettledAmount         float64     `json:"settled_amount"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	DeletedAt             nulls.Time  `json:"deleted_at"`
	ProgressEventID       nulls.Int32 `json:"progress_event_id"`
	TransactionStatusID   int16       `json:"transaction_status_id"`
	StatusName            string      `json:"status_name"`
	StatusDescription     string      `json:"status_description"`
	TransactionTypeID     int16       `json:"transaction_type_id"`
	TypeName              string      `json:"type_name"`
	TypeDescription       string      `json:"type_description"`
}

func (q *Queries) AllMerchantTrx(ctx context.Context, arg AllMerchantTrxParams) ([]AllMerchantTrxRow, error) {
	rows, err := q.db.QueryContext(ctx, allMerchantTrx, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AllMerchantTrxRow{}
	for rows.Next() {
		var i AllMerchantTrxRow
		if err := rows.Scan(
			&i.ID,
			&i.MerchantTransactionID,
			&i.OwnerID,
			&i.TransactionID,
			&i.TransactionDate,
			&i.TransactionDatetime,
			&i.CollectedAmount,
			&i.SettledAmount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ProgressEventID,
			&i.TransactionStatusID,
			&i.StatusName,
			&i.StatusDescription,
			&i.TransactionTypeID,
			&i.TypeName,
			&i.TypeDescription,
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

const createMerchantTrx = `-- name: CreateMerchantTrx :one
INSERT INTO merchant_transactions (
    transaction_status_id, transaction_type_id, progress_event_id,
    merchant_transaction_id, owner_id, transaction_id,
    transaction_date, transaction_datetime, collected_amount,
    settled_amount, created_at, updated_at, deleted_at
) VALUES (
    $1, $2, NULL, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING id, transaction_status_id, transaction_type_id, progress_event_id, merchant_transaction_id, owner_id, transaction_id, transaction_date, transaction_datetime, collected_amount, settled_amount, created_at, updated_at, deleted_at
`

type CreateMerchantTrxParams struct {
	TransactionStatusID   int16      `json:"transaction_status_id"`
	TransactionTypeID     int16      `json:"transaction_type_id"`
	MerchantTransactionID string     `json:"merchant_transaction_id"`
	OwnerID               string     `json:"owner_id"`
	TransactionID         string     `json:"transaction_id"`
	TransactionDate       time.Time  `json:"transaction_date"`
	TransactionDatetime   time.Time  `json:"transaction_datetime"`
	CollectedAmount       float64    `json:"collected_amount"`
	SettledAmount         float64    `json:"settled_amount"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             nulls.Time `json:"deleted_at"`
}

func (q *Queries) CreateMerchantTrx(ctx context.Context, arg CreateMerchantTrxParams) (MerchantTransaction, error) {
	row := q.db.QueryRowContext(ctx, createMerchantTrx,
		arg.TransactionStatusID,
		arg.TransactionTypeID,
		arg.MerchantTransactionID,
		arg.OwnerID,
		arg.TransactionID,
		arg.TransactionDate,
		arg.TransactionDatetime,
		arg.CollectedAmount,
		arg.SettledAmount,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.DeletedAt,
	)
	var i MerchantTransaction
	err := row.Scan(
		&i.ID,
		&i.TransactionStatusID,
		&i.TransactionTypeID,
		&i.ProgressEventID,
		&i.MerchantTransactionID,
		&i.OwnerID,
		&i.TransactionID,
		&i.TransactionDate,
		&i.TransactionDatetime,
		&i.CollectedAmount,
		&i.SettledAmount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteDuplicateMerchantTrx = `-- name: DeleteDuplicateMerchantTrx :one
WITH deleted AS (
	DELETE FROM merchant_transactions
	WHERE id IN (
		SELECT id from (
		  SELECT id,
		  ROW_NUMBER() OVER(PARTITION BY merchant_transaction_id, transaction_date ORDER BY id ASC) AS result
		  FROM merchant_transactions
		  WHERE merchant_transactions.transaction_date BETWEEN $1::date AND $2::date
          AND merchant_transactions.owner_id = $3::text
		) duplicates
		WHERE duplicates.result > 1
	)
	RETURNING id
) SELECT count(id) FROM deleted
`

type DeleteDuplicateMerchantTrxParams struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	PlatformID string    `json:"platform_id"`
}

func (q *Queries) DeleteDuplicateMerchantTrx(ctx context.Context, arg DeleteDuplicateMerchantTrxParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteDuplicateMerchantTrx, arg.StartDate, arg.EndDate, arg.PlatformID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteMerchantTrxByID = `-- name: DeleteMerchantTrxByID :exec
DELETE FROM merchant_transactions
WHERE id = $1
`

func (q *Queries) DeleteMerchantTrxByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMerchantTrxByID, id)
	return err
}