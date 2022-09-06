-- name: CreateMerchantTrx :one
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
RETURNING *;

-- name: AllMerchantTrx :many
SELECT 
    p.id, p.merchant_transaction_id, p.merchant_transaction_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM merchant_transactions p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
OFFSET $1
LIMIT $2;

-- name: AllDuplicateMerchantTrx :many 
SELECT 
    p.id, p.merchant_transaction_id, p.merchant_transaction_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM (
    SELECT *,
    ROW_NUMBER() OVER(PARTITION BY merchant_transaction_id, transaction_date ORDER BY id asc) AS res
    FROM merchant_transactions
) p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
where p.res > 1
OFFSET $1
LIMIT $2;

-- name: DeleteDuplicateMerchantTrx :one
WITH deleted AS (
	DELETE FROM merchant_transactions
	WHERE id IN (
		SELECT id from (
		  SELECT id,
		  ROW_NUMBER() OVER(PARTITION BY merchant_transaction_id, transaction_date ORDER BY id ASC) AS result
		  FROM merchant_transactions
		  WHERE merchant_transactions.transaction_date BETWEEN sqlc.arg(start_date)::date AND sqlc.arg(end_date)::date
          AND merchant_transactions.owner_id = sqlc.arg(platform_id)::int
		) duplicates
		WHERE duplicates.result > 1
	)
	RETURNING id
) SELECT count(id) FROM deleted;

-- name: DeleteMerchantTrxByID :exec
DELETE FROM merchant_transactions
WHERE id = $1;