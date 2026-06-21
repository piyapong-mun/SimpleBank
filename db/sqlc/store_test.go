package db

import (
	"context"
	"fmt"
	"testing"

	// import uuid

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(TestDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	startAccount1Balance := account1.Balance
	startAccount2Balance := account2.Balance

	numCalls := 5
	amount := int64(10)

	fmt.Println("Before", account1.Balance, account2.Balance)

	results := make(chan TransferTxResult, numCalls)
	errChan := make(chan error, numCalls)

	for i := 0; i < numCalls; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errChan <- err
			results <- result
		}()
	}

	var resultsList []TransferTxResult
	for i := 0; i < numCalls; i++ {
		result := <-results
		err := <-errChan
		assert.NoError(t, err)
		assert.Equal(t, amount, result.Transfer.Amount)

		// check result
		assert.NotZero(t, result.Transfer.ID)
		assert.Equal(t, result.FromAccount.ID, account1.ID)
		assert.Equal(t, result.ToAccount.ID, account2.ID)
		assert.Equal(t, amount, result.Transfer.Amount)
		assert.Equal(t, result.FromAccount.ID, result.FromEntry.AccountID)
		assert.Equal(t, result.ToAccount.ID, result.ToEntry.AccountID)
		assert.Equal(t, amount, result.FromEntry.Amount)
		assert.Equal(t, -amount, result.ToEntry.Amount)

		// check account balance
		nowbalance1 := startAccount1Balance - amount*int64(i+1)
		nowbalance2 := startAccount2Balance + amount*int64(i+1)
		assert.Equal(t, result.FromAccount.Balance, nowbalance1)
		assert.Equal(t, result.ToAccount.Balance, nowbalance2)
		fmt.Println("After", result.FromAccount.Balance, result.ToAccount.Balance)

		resultsList = append(resultsList, result)
	}

	// Cleanup
	for _, result := range resultsList {
		err := TestQueries.DeleteTransfers(context.Background(), result.Transfer.ID)
		assert.NoError(t, err)

		err = TestQueries.DeleteEntry(context.Background(), result.FromEntry.ID)
		assert.NoError(t, err)

		err = TestQueries.DeleteEntry(context.Background(), result.ToEntry.ID)
		assert.NoError(t, err)
	}

	err := TestQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	err = TestQueries.DeleteAccount(context.Background(), account2.ID)
	assert.NoError(t, err)
}

func TestDeadLock(t *testing.T) {
	store := NewStore(TestDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	numCalls := 10
	amount := int64(10)

	fmt.Println("Before", account1.Balance, account2.Balance)

	errChan := make(chan error, numCalls)
	results := make(chan TransferTxResult, numCalls)

	for i := 0; i < numCalls; i++ {

		// acc1 --> acc2
		FromAccountID := account1.ID
		ToAccountID := account2.ID
		if i%2 == 0 {
			// acc2 --> acc1
			FromAccountID = account2.ID
			ToAccountID = account1.ID
		}

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})

			errChan <- err
			results <- result
		}()
	}

	var resultsList []TransferTxResult
	for i := 0; i < numCalls; i++ {
		err := <-errChan
		assert.NoError(t, err)
		result := <-results
		resultsList = append(resultsList, result)
	}

	// Cleanup
	for _, result := range resultsList {
		err := TestQueries.DeleteTransfers(context.Background(), result.Transfer.ID)
		assert.NoError(t, err)

		err = TestQueries.DeleteEntry(context.Background(), result.FromEntry.ID)
		assert.NoError(t, err)

		err = TestQueries.DeleteEntry(context.Background(), result.ToEntry.ID)
		assert.NoError(t, err)
	}

	err := TestQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	err = TestQueries.DeleteAccount(context.Background(), account2.ID)
	assert.NoError(t, err)
}
