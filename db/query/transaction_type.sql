-- name: CreateTransactionType :one
INSERT INTO transaction_types (
    type_name, type_description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateTransactionType :one
UPDATE transaction_types 
SET 
    type_name = CASE WHEN sqlc.arg(type_name)::text <> '' THEN sqlc.arg(type_name)::text ELSE type_name END,
    type_description = CASE WHEN sqlc.arg(type_description)::text <> '' THEN sqlc.arg(type_description)::text ELSE type_description END,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTransactionType :exec
DELETE FROM transaction_types 
WHERE id = $1;

-- name: ViewTransactionType :one
SELECT *
FROM transaction_types
WHERE id = $1 LIMIT 1;

-- name: AllTransactionType :many
SELECT *
FROM transaction_types
ORDER BY created_at
OFFSET $1
LIMIT $2;