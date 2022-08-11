-- name: CreateProgressEventType :one
INSERT INTO progress_event_types (
    progress_event_type_name, progress_event_type_description
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateProgressEventType :one
UPDATE progress_event_types 
SET 
    progress_event_type_name = CASE WHEN sqlc.arg(progress_event_type_name)::text <> '' THEN sqlc.arg(progress_event_type_name)::text ELSE progress_event_type_name END,
    progress_event_type_description = CASE WHEN sqlc.arg(progress_event_type_description)::text <> '' THEN sqlc.arg(progress_event_type_description)::text ELSE progress_event_type_description END,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProgressEventType :exec
DELETE FROM progress_event_types 
WHERE id = $1;

-- name: ViewProgressEventType :one
SELECT *
FROM progress_event_types
WHERE id = $1 LIMIT 1;

-- name: AllProgressEventType :many
SELECT *
FROM progress_event_types
ORDER BY created_at
OFFSET $1
LIMIT $2;