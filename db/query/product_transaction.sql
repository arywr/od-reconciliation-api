-- name: CreateProductTransaction :one
INSERT INTO product_transactions (
    transaction_status_id, transaction_type_id, progress_event_id,
    product_transaction_id, merchant_transaction_id, channel_transaction_id,
    owner_id, transaction_id, transaction_date, 
    transaction_datetime, collected_amount, settled_amount, 
    created_at, updated_at, deleted_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: AllProductTransaction :many
SELECT
    p.id, p.product_transaction_id, p.merchant_transaction_id, p.channel_transaction_id, p.owner_id, p.transaction_id, p.transaction_date, p.transaction_datetime, p.collected_amount, p.settled_amount, p.created_at, p.updated_at, p.deleted_at,
    p.progress_event_id, q.progress_event_type_id, q.progress_name, q.status, q.percentage, q.file,
    p.transaction_status_id, s.status_name, s.status_description,
    p.transaction_type_id, r.type_name, r.type_description,
    t.progress_event_type_name, t.progress_event_type_description
FROM product_transactions AS p
JOIN progress_events AS q ON p.progress_event_id = q.id
JOIN transaction_types AS r ON p.transaction_type_id = r.id
JOIN transaction_statuses AS s ON p.transaction_status_id = s.id
JOIN progress_event_types AS t ON q.progress_event_type_id = t.id;