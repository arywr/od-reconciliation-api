-- name: GetReconcileRule :many
SELECT *
FROM reconcile_rules
WHERE product_id = $1
AND platform_id = $2;