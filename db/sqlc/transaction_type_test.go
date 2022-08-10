package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/arywr/od-reconciliation-api/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomType(t *testing.T) OdTransactionType {
	args := CreateTransactionTypeParams{
		TypeName:        util.RandomType(),
		TypeDescription: util.RandomDescription(),
	}

	types, err := testQueries.CreateTransactionType(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, types)

	require.Equal(t, args.TypeName, types.TypeName)
	require.Equal(t, args.TypeDescription, types.TypeDescription)

	require.NotZero(t, types.ID)
	require.NotZero(t, types.CreatedAt)
	require.NotZero(t, types.UpdatedAt)

	return types
}

func TestCreateTransactionType(t *testing.T) {
	CreateRandomType(t)
}

func TestGetTransactionType(t *testing.T) {
	type1 := CreateRandomType(t)
	type2, err := testQueries.ViewTransactionType(context.Background(), type1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, type2)

	require.Equal(t, type1.ID, type2.ID)
	require.Equal(t, type1.TypeName, type2.TypeName)
	require.Equal(t, type1.TypeDescription, type2.TypeDescription)
	require.WithinDuration(t, type1.CreatedAt, type2.CreatedAt, time.Second)
}

func TestUpdateTransactionType(t *testing.T) {
	type1 := CreateRandomType(t)

	args := UpdateTransactionTypeParams{
		ID:       type1.ID,
		TypeName: util.RandomType(),
	}

	type2, err := testQueries.UpdateTransactionType(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, type2)

	require.Equal(t, type1.ID, type2.ID)
	require.Equal(t, args.TypeName, type2.TypeName)
}

func TestDeleteTransactionType(t *testing.T) {
	type1 := CreateRandomType(t)
	err := testQueries.DeleteTransactionType(context.Background(), type1.ID)
	require.NoError(t, err)

	type2, err := testQueries.ViewTransactionType(context.Background(), type1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, type2)
}

func TestAllTransactionType(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomType(t)
	}

	args := AllTransactionTypeParams{
		Limit:  10,
		Offset: 0,
	}

	types, err := testQueries.AllTransactionType(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, types, 10)

	for _, row := range types {
		require.NotEmpty(t, row)
	}
}
