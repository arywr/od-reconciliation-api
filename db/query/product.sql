-- name: GetRelatedPlatforms :many
SELECT products.id product_id, products.product_name, products.product_has_sub, sub_products.id sub_product_id, sub_product_name, merchants.id merchant_id, merchants.merchant_name merchant_name
FROM products
LEFT OUTER JOIN sub_products ON sub_products.product_id = products.id
LEFT OUTER JOIN merchants ON merchants.product_id = products.id
WHERE products.id = $1;