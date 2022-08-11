-- name: CreateTransactionStatus :one
INSERT INTO transaction_statuses (
    status_name, status_description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateTransactionStatus :one
UPDATE transaction_statuses 
SET 
    status_name = CASE WHEN sqlc.arg(status_name)::text <> '' THEN sqlc.arg(status_name)::text ELSE status_name END,
    status_description = CASE WHEN sqlc.arg(status_description)::text <> '' THEN sqlc.arg(status_description)::text ELSE status_description END,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTransactionStatus :exec
DELETE FROM transaction_statuses
WHERE id = $1;

-- name: ViewTransactionStatus :one
SELECT *
FROM transaction_statuses
WHERE id = $1 LIMIT 1;

-- name: AllTransactionStatus :many
SELECT *
FROM transaction_statuses
ORDER BY created_at
OFFSET $1
LIMIT $2;