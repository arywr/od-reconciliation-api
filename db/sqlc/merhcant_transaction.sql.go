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
    p.id, p.merchant_transaction_id, p.merchant_transaction_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM (
    SELECT id, product_id, sub_product_id, platform_id, sub_platform_id, transaction_status_id, transaction_type_id, progress_event_id, product_transaction_id, merchant_transaction_id, transaction_id, transaction_date, transaction_datetime, channel_code, channel_name, merchant_code, merchant_name, product_code, product_name, collected_amount, settled_amount, reconcile_at, created_at, updated_at, deleted_at,
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
	ID                      int64        `json:"id"`
	MerchantTransactionID   nulls.String `json:"merchant_transaction_id"`
	MerchantTransactionID_2 nulls.String `json:"merchant_transaction_id_2"`
	TransactionID           string       `json:"transaction_id"`
	TransactionDate         time.Time    `json:"transaction_date"`
	TransactionDatetime     time.Time    `json:"transaction_datetime"`
	CollectedAmount         float64      `json:"collected_amount"`
	SettledAmount           float64      `json:"settled_amount"`
	CreatedAt               time.Time    `json:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at"`
	DeletedAt               nulls.Time   `json:"deleted_at"`
	ProgressEventID         nulls.Int32  `json:"progress_event_id"`
	TransactionStatusID     int16        `json:"transaction_status_id"`
	StatusName              string       `json:"status_name"`
	StatusDescription       string       `json:"status_description"`
	TransactionTypeID       int16        `json:"transaction_type_id"`
	TypeName                string       `json:"type_name"`
	TypeDescription         string       `json:"type_description"`
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
			&i.MerchantTransactionID_2,
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
    p.id, p.merchant_transaction_id, p.merchant_transaction_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
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
	ID                      int64        `json:"id"`
	MerchantTransactionID   nulls.String `json:"merchant_transaction_id"`
	MerchantTransactionID_2 nulls.String `json:"merchant_transaction_id_2"`
	TransactionID           string       `json:"transaction_id"`
	TransactionDate         time.Time    `json:"transaction_date"`
	TransactionDatetime     time.Time    `json:"transaction_datetime"`
	CollectedAmount         float64      `json:"collected_amount"`
	SettledAmount           float64      `json:"settled_amount"`
	CreatedAt               time.Time    `json:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at"`
	DeletedAt               nulls.Time   `json:"deleted_at"`
	ProgressEventID         nulls.Int32  `json:"progress_event_id"`
	TransactionStatusID     int16        `json:"transaction_status_id"`
	StatusName              string       `json:"status_name"`
	StatusDescription       string       `json:"status_description"`
	TransactionTypeID       int16        `json:"transaction_type_id"`
	TypeName                string       `json:"type_name"`
	TypeDescription         string       `json:"type_description"`
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
			&i.MerchantTransactionID_2,
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
    product_transaction_id, merchant_transaction_id, product_id, sub_product_id,
    platform_id, sub_platform_id, transaction_id, transaction_date, transaction_datetime, 
    channel_code, channel_name, merchant_code, merchant_name,
    product_code, product_name, collected_amount, settled_amount, 
    created_at, updated_at, deleted_at
) VALUES (
    $1, $2, NULL, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
)
RETURNING id, product_id, sub_product_id, platform_id, sub_platform_id, transaction_status_id, transaction_type_id, progress_event_id, product_transaction_id, merchant_transaction_id, transaction_id, transaction_date, transaction_datetime, channel_code, channel_name, merchant_code, merchant_name, product_code, product_name, collected_amount, settled_amount, reconcile_at, created_at, updated_at, deleted_at
`

type CreateMerchantTrxParams struct {
	TransactionStatusID   int16        `json:"transaction_status_id"`
	TransactionTypeID     int16        `json:"transaction_type_id"`
	ProductTransactionID  nulls.String `json:"product_transaction_id"`
	MerchantTransactionID nulls.String `json:"merchant_transaction_id"`
	ProductID             int32        `json:"product_id"`
	SubProductID          nulls.Int32  `json:"sub_product_id"`
	PlatformID            nulls.Int32  `json:"platform_id"`
	SubPlatformID         nulls.Int32  `json:"sub_platform_id"`
	TransactionID         string       `json:"transaction_id"`
	TransactionDate       time.Time    `json:"transaction_date"`
	TransactionDatetime   time.Time    `json:"transaction_datetime"`
	ChannelCode           nulls.String `json:"channel_code"`
	ChannelName           nulls.String `json:"channel_name"`
	MerchantCode          nulls.String `json:"merchant_code"`
	MerchantName          nulls.String `json:"merchant_name"`
	ProductCode           nulls.String `json:"product_code"`
	ProductName           nulls.String `json:"product_name"`
	CollectedAmount       float64      `json:"collected_amount"`
	SettledAmount         float64      `json:"settled_amount"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
	DeletedAt             nulls.Time   `json:"deleted_at"`
}

func (q *Queries) CreateMerchantTrx(ctx context.Context, arg CreateMerchantTrxParams) (MerchantTransaction, error) {
	row := q.db.QueryRowContext(ctx, createMerchantTrx,
		arg.TransactionStatusID,
		arg.TransactionTypeID,
		arg.ProductTransactionID,
		arg.MerchantTransactionID,
		arg.ProductID,
		arg.SubProductID,
		arg.PlatformID,
		arg.SubPlatformID,
		arg.TransactionID,
		arg.TransactionDate,
		arg.TransactionDatetime,
		arg.ChannelCode,
		arg.ChannelName,
		arg.MerchantCode,
		arg.MerchantName,
		arg.ProductCode,
		arg.ProductName,
		arg.CollectedAmount,
		arg.SettledAmount,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.DeletedAt,
	)
	var i MerchantTransaction
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.SubProductID,
		&i.PlatformID,
		&i.SubPlatformID,
		&i.TransactionStatusID,
		&i.TransactionTypeID,
		&i.ProgressEventID,
		&i.ProductTransactionID,
		&i.MerchantTransactionID,
		&i.TransactionID,
		&i.TransactionDate,
		&i.TransactionDatetime,
		&i.ChannelCode,
		&i.ChannelName,
		&i.MerchantCode,
		&i.MerchantName,
		&i.ProductCode,
		&i.ProductName,
		&i.CollectedAmount,
		&i.SettledAmount,
		&i.ReconcileAt,
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
          AND merchant_transactions.owner_id = $3::int
		) duplicates
		WHERE duplicates.result > 1
	)
	RETURNING id
) SELECT count(id) FROM deleted
`

type DeleteDuplicateMerchantTrxParams struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	PlatformID int32     `json:"platform_id"`
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
