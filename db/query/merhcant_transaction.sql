-- name: CreateMerchantTrx :one
INSERT INTO merchant_transactions (
    transaction_status_id, transaction_type_id, progress_event_id,
    merchant_transaction_id, owner_id, transaction_id,
    transaction_date, transaction_datetime, collected_amount,
    settled_amount, created_at, updated_at, deleted_at
) VALUES (
    $1, $2, NULL, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: AllMerchantTrx :many
SELECT 
    p.id, p.merchant_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM merchant_transactions p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
OFFSET $1
LIMIT $2;

-- name: AllDuplicateMerchantTrx :many 
SELECT 
    p.id, p.merchant_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
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
          AND merchant_transactions.owner_id = sqlc.arg(platform_id)::text
		) duplicates
		WHERE duplicates.result > 1
	)
	RETURNING id
) SELECT count(id) FROM deleted;

-- name: DeleteMerchantTrxByID :exec
DELETE FROM merchant_transactions
WHERE id = $1;