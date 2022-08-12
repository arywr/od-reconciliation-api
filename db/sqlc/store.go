package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction Error %s %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type ProgressEventTypeFromParent struct {
	ID                           int64  `json:"id"`
	ProgressEventTypeName        string `json:"progress_event_type_name"`
	ProgressEventTypeDescription string `json:"progress_event_type_description"`
}

type UpdateProgressTxParams struct {
	ID         int64   `json:"progress_event_id"`
	Percentage float64 `json:"percentage" binding:"min=0"`
	Status     string  `json:"status"`
}

type UpdateProgressTxResult struct {
	ID                int64                       `json:"id"`
	ProgressEventType ProgressEventTypeFromParent `json:"progress_event_type"`
	ProgressName      string                      `json:"progress_name"`
	Status            string                      `json:"status"`
	Percentage        float64                     `json:"percentage"`
	File              string                      `json:"file"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
	DeletedAt         *time.Time                  `json:"deleted_at"`
}

func (store *Store) UpdateProgressTx(ctx context.Context, arg UpdateProgressTxParams) (UpdateProgressTxResult, error) {
	var result UpdateProgressTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		args := UpdateProgressParams{
			ID:         arg.ID,
			Percentage: arg.Percentage,
			Status:     arg.Status,
		}

		progress, err := q.UpdateProgress(ctx, args)
		if err != nil {
			return err
		}

		evtType, err := q.ViewProgressEventType(ctx, int64(progress.ProgressEventTypeID))
		if err != nil {
			return err
		}

		result.ID = progress.ID
		result.ProgressName = progress.ProgressName
		result.ProgressEventType.ID = evtType.ID
		result.ProgressEventType.ProgressEventTypeName = evtType.ProgressEventTypeName
		result.ProgressEventType.ProgressEventTypeDescription = evtType.ProgressEventTypeDescription
		result.Status = progress.Status
		result.Percentage = progress.Percentage
		result.File = progress.File
		result.CreatedAt = progress.CreatedAt
		result.UpdatedAt = progress.UpdatedAt
		result.DeletedAt = progress.DeletedAt

		return nil
	})

	return result, err
}
