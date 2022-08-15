-- name: MatchReconciliationProduct :exec
WITH join_product AS (
    SELECT t1.id
    FROM product_transactions t1
    LEFT OUTER JOIN merchant_transactions t2 ON t1.product_transaction_id = t2.merchant_transaction_id
    WHERE t2.merchant_transaction_id IS NOT NULL
    AND t1.owner_id = sqlc.arg(platform_id)
    AND t1.created_at BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
    AND t2.created_at BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
)
UPDATE product_transactions
SET transaction_status_id = 3, updated_at = now()
FROM join_product
WHERE product_transactions.id = join_product.id;

-- name: MatchReconciliationMerchant :exec
WITH join_merchant AS (
    SELECT t1.id
    FROM merchant_transactions t1
    LEFT OUTER JOIN product_transactions t2 ON t1.merchant_transaction_id = t2.product_transaction_id
    WHERE t2.product_transaction_id IS NOT NULL
    AND t1.owner_id = sqlc.arg(destination_id)
    AND t1.created_at BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
    AND t2.created_at BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
)
UPDATE merchant_transactions
SET transaction_status_id = 3, updated_at = now()
FROM join_merchant
WHERE merchant_transactions.id = join_merchant.id;

-- name: DisputeReconciliationProduct :exec
WITH join_product_dispute AS (
    SELECT t1.id
    FROM product_transactions t1
    WHERE t1.transaction_status_id != 2
    AND t1.transaction_status_id != 3
    AND t1.owner_id = sqlc.arg(platform_id)
    AND t1.created_at BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
)
UPDATE product_transactions
SET transaction_status_id = 4, updated_at = now()
FROM join_product_dispute
WHERE product_transactions.id = join_product_dispute.id;

-- name: DisputeReconciliationMerchant :exec
WITH join_merchant_dispute AS (
    SELECT t1.id
    FROM merchant_transactions t1
    WHERE t1.transaction_status_id != 2
    AND t1.transaction_status_id != 3
    AND t1.owner_id = sqlc.arg(destination_id)
    AND t1.created_at BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
)
UPDATE merchant_transactions
SET transaction_status_id = 4, updated_at = now()
FROM join_merchant_dispute
WHERE merchant_transactions.id = join_merchant_dispute.id;