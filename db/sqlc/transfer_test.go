package db

import (
	"context"
	"testing"
	"time"

	"github.com/piyapong-mun/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
	arg := CreateTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomNumber(1, 1000),
	}

	transfer, err := TestQueries.CreateTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	assert.Equal(t, arg.Amount, transfer.Amount)
	assert.NotZero(t, transfer.ID)
	assert.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	transfer := createRandomTransfer(t, fromAccount, toAccount)

	// Cleanup
	err := TestQueries.DeleteTransfers(context.Background(), transfer.ID)
	require.NoError(t, err)
	clearAccount(t, fromAccount.ID)
	clearAccount(t, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, fromAccount, toAccount)

	transfer2, err := TestQueries.GetTransfers(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	assert.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

	// Cleanup
	err = TestQueries.DeleteTransfers(context.Background(), transfer1.ID)
	require.NoError(t, err)
	clearAccount(t, fromAccount.ID)
	clearAccount(t, toAccount.ID)
}

func TestUpdateTransfer(t *testing.T) {
	fromAccount1 := createRandomAccount(t)
	toAccount1 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, fromAccount1, toAccount1)

	fromAccount2 := createRandomAccount(t)
	toAccount2 := createRandomAccount(t)

	arg := UpdateTransfersParams{
		ID:            transfer1.ID,
		FromAccountID: fromAccount2.ID,
		ToAccountID:   toAccount2.ID,
		Amount:        util.RandomNumber(1, 1000),
	}

	err := TestQueries.UpdateTransfers(context.Background(), arg)
	require.NoError(t, err)

	transfer2, err := TestQueries.GetTransfers(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, arg.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, arg.Amount, transfer2.Amount)
	assert.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

	// Cleanup
	err = TestQueries.DeleteTransfers(context.Background(), transfer1.ID)
	require.NoError(t, err)
	clearAccount(t, fromAccount1.ID)
	clearAccount(t, toAccount1.ID)
	clearAccount(t, fromAccount2.ID)
	clearAccount(t, toAccount2.ID)
}

func TestDeleteTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, fromAccount, toAccount)

	err := TestQueries.DeleteTransfers(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := TestQueries.GetTransfers(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Empty(t, transfer2)

	// Cleanup
	clearAccount(t, fromAccount.ID)
	clearAccount(t, toAccount.ID)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	var lastTransfers []Transfer
	for i := 0; i < 10; i++ {
		lastTransfers = append(lastTransfers, createRandomTransfer(t, fromAccount, toAccount))
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := TestQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

	// Cleanup
	for _, transfer := range lastTransfers {
		err = TestQueries.DeleteTransfers(context.Background(), transfer.ID)
		require.NoError(t, err)
	}
	clearAccount(t, fromAccount.ID)
	clearAccount(t, toAccount.ID)
}
