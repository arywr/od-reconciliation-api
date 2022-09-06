-- name: GetTransactionModel :one
SELECT transaction_models.*
FROM products
JOIN transaction_models ON products.transaction_model_id = transaction_models.id
WHERE products.id = sqlc.arg(platform_id)::int;

-- name: GetTransactionModelMerchant :one
SELECT transaction_models.*
FROM merchants
JOIN transaction_models ON merchants.transaction_model_id = transaction_models.id
WHERE merchants.id = sqlc.arg(platform_id)::int;