-- name: CreateTransactionType :one
INSERT INTO od_transaction_types (
    type_name, type_description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateTransactionType :one
UPDATE od_transaction_types 
SET type_name = $1, type_description = $2
WHERE id = $3
RETURNING *;

-- name: DeleteTransactionType :exec
DELETE FROM od_transaction_types 
WHERE id = $1;

-- name: ViewTransactionType :one
SELECT id, type_name, type_description
FROM od_transaction_types
WHERE id = $1 LIMIT 1;

-- name: AllTransactionType :many
SELECT id, type_name, type_description
FROM od_transaction_types
ORDER BY created_at;