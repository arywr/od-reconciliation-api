-- name: CreateProductTransaction :one
INSERT INTO product_transactions (
    transaction_status_id, transaction_type_id, progress_event_id,
    product_transaction_id, merchant_transaction_id, channel_transaction_id,
    owner_id, transaction_id, transaction_date, 
    transaction_datetime, collected_amount, settled_amount, 
    created_at, updated_at, deleted_at
) VALUES (
    $1, $2, NULL, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: AllProductTransaction :many
SELECT
    p.id, p.product_transaction_id, p.merchant_transaction_id, p.channel_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description
FROM product_transactions AS p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
OFFSET $1
LIMIT $2;

-- name: AllDuplicateProductTransaction :many
SELECT p.id, p.product_transaction_id, p.merchant_transaction_id, p.channel_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at, p.progress_event_id,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description from (
  SELECT *,
  ROW_NUMBER() OVER(PARTITION BY product_transaction_id, transaction_date ORDER BY id asc) AS res
  FROM product_transactions
) p
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
where p.res > 1
OFFSET $1
LIMIT $2;

-- name: DeleteDuplicateProductTrx :one
WITH deleted AS (
	DELETE FROM product_transactions
	WHERE id IN (
		SELECT id from (
		  SELECT id,
		  ROW_NUMBER() OVER(PARTITION BY product_transaction_id, transaction_date ORDER BY id ASC) AS result
		  FROM product_transactions
		  WHERE product_transactions.transaction_date BETWEEN sqlc.arg(start_date)::date AND sqlc.arg(end_date)::date
          AND product_transactions.owner_id = sqlc.arg(platform_id)::text
		) duplicates
		WHERE duplicates.result > 1
	)
	RETURNING id
) SELECT count(id) FROM deleted;

-- name: DeleteProductTrxByID :exec
DELETE FROM product_transactions
WHERE id = $1;