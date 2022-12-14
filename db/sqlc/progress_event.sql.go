// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: progress_event.sql

package db

import (
	"context"
	"time"

	"github.com/gobuffalo/nulls"
)

const allProgressEvent = `-- name: AllProgressEvent :many
SELECT 
    p.id, p.progress_event_type_id ,p.progress_name, p.status, p.percentage, p.file, p.created_at, p.updated_at, p.deleted_at,
    q.progress_event_type_name, q.progress_event_type_description
FROM progress_events AS p
JOIN progress_event_types AS q ON p.progress_event_type_id = q.id
ORDER BY p.created_at ASC
OFFSET $1
LIMIT $2
`

type AllProgressEventParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type AllProgressEventRow struct {
	ID                           int64      `json:"id"`
	ProgressEventTypeID          int16      `json:"progress_event_type_id"`
	ProgressName                 string     `json:"progress_name"`
	Status                       string     `json:"status"`
	Percentage                   float64    `json:"percentage"`
	File                         string     `json:"file"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`
	DeletedAt                    nulls.Time `json:"deleted_at"`
	ProgressEventTypeName        string     `json:"progress_event_type_name"`
	ProgressEventTypeDescription string     `json:"progress_event_type_description"`
}

func (q *Queries) AllProgressEvent(ctx context.Context, arg AllProgressEventParams) ([]AllProgressEventRow, error) {
	rows, err := q.db.QueryContext(ctx, allProgressEvent, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AllProgressEventRow{}
	for rows.Next() {
		var i AllProgressEventRow
		if err := rows.Scan(
			&i.ID,
			&i.ProgressEventTypeID,
			&i.ProgressName,
			&i.Status,
			&i.Percentage,
			&i.File,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ProgressEventTypeName,
			&i.ProgressEventTypeDescription,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createProgressEvent = `-- name: CreateProgressEvent :one
INSERT INTO progress_events (
    progress_event_type_id, progress_name, status, percentage, file
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, progress_event_type_id, progress_name, status, percentage, file, created_at, updated_at, deleted_at
`

type CreateProgressEventParams struct {
	ProgressEventTypeID int16   `json:"progress_event_type_id"`
	ProgressName        string  `json:"progress_name"`
	Status              string  `json:"status"`
	Percentage          float64 `json:"percentage"`
	File                string  `json:"file"`
}

func (q *Queries) CreateProgressEvent(ctx context.Context, arg CreateProgressEventParams) (ProgressEvent, error) {
	row := q.db.QueryRowContext(ctx, createProgressEvent,
		arg.ProgressEventTypeID,
		arg.ProgressName,
		arg.Status,
		arg.Percentage,
		arg.File,
	)
	var i ProgressEvent
	err := row.Scan(
		&i.ID,
		&i.ProgressEventTypeID,
		&i.ProgressName,
		&i.Status,
		&i.Percentage,
		&i.File,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteProgressEvent = `-- name: DeleteProgressEvent :exec
DELETE FROM progress_events
WHERE id = $1
`

func (q *Queries) DeleteProgressEvent(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProgressEvent, id)
	return err
}

const updateProgress = `-- name: UpdateProgress :one
UPDATE progress_events
SET
    percentage = $1,
    status = CASE WHEN $3::text <> '' THEN $3::text ELSE status END,
    updated_at = now()
WHERE id = $2
RETURNING id, progress_event_type_id, progress_name, status, percentage, file, created_at, updated_at, deleted_at
`

type UpdateProgressParams struct {
	Percentage float64 `json:"percentage"`
	ID         int64   `json:"id"`
	Status     string  `json:"status"`
}

func (q *Queries) UpdateProgress(ctx context.Context, arg UpdateProgressParams) (ProgressEvent, error) {
	row := q.db.QueryRowContext(ctx, updateProgress, arg.Percentage, arg.ID, arg.Status)
	var i ProgressEvent
	err := row.Scan(
		&i.ID,
		&i.ProgressEventTypeID,
		&i.ProgressName,
		&i.Status,
		&i.Percentage,
		&i.File,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const viewProgressEvent = `-- name: ViewProgressEvent :one
SELECT 
    p.id, p.progress_event_type_id ,p.progress_name, p.status, p.percentage, p.file, p.created_at, p.updated_at, p.deleted_at,
    q.progress_event_type_name, q.progress_event_type_description
FROM progress_events AS p
JOIN progress_event_types AS q ON p.progress_event_type_id = q.id
WHERE p.id = $1 LIMIT 1
`

type ViewProgressEventRow struct {
	ID                           int64      `json:"id"`
	ProgressEventTypeID          int16      `json:"progress_event_type_id"`
	ProgressName                 string     `json:"progress_name"`
	Status                       string     `json:"status"`
	Percentage                   float64    `json:"percentage"`
	File                         string     `json:"file"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`
	DeletedAt                    nulls.Time `json:"deleted_at"`
	ProgressEventTypeName        string     `json:"progress_event_type_name"`
	ProgressEventTypeDescription string     `json:"progress_event_type_description"`
}

func (q *Queries) ViewProgressEvent(ctx context.Context, id int64) (ViewProgressEventRow, error) {
	row := q.db.QueryRowContext(ctx, viewProgressEvent, id)
	var i ViewProgressEventRow
	err := row.Scan(
		&i.ID,
		&i.ProgressEventTypeID,
		&i.ProgressName,
		&i.Status,
		&i.Percentage,
		&i.File,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.ProgressEventTypeName,
		&i.ProgressEventTypeDescription,
	)
	return i, err
}
