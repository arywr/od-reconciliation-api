-- name: CreateProgressEvent :one
INSERT INTO progress_events (
    progress_event_type_id, progress_name, status, percentage, file
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateProgress :one
UPDATE progress_events
SET
    percentage = $1,
    status = CASE WHEN sqlc.arg(status)::text <> '' THEN sqlc.arg(status)::text ELSE status END,
    updated_at = now()
WHERE id = $2
RETURNING *;

-- name: DeleteProgressEvent :exec
DELETE FROM progress_events
WHERE id = $1;

-- name: ViewProgressEvent :one
SELECT 
    p.id, p.progress_event_type_id ,p.progress_name, p.status, p.percentage, p.file, p.created_at, p.updated_at, p.deleted_at,
    q.progress_event_type_name, q.progress_event_type_description
FROM progress_events AS p
JOIN progress_event_types AS q ON p.progress_event_type_id = q.id
WHERE p.id = $1 LIMIT 1;

-- name: AllProgressEvent :many
SELECT 
    p.id, p.progress_event_type_id ,p.progress_name, p.status, p.percentage, p.file, p.created_at, p.updated_at, p.deleted_at,
    q.progress_event_type_name, q.progress_event_type_description
FROM progress_events AS p
JOIN progress_event_types AS q ON p.progress_event_type_id = q.id
ORDER BY p.created_at ASC
OFFSET $1
LIMIT $2;