package db

import (
	"context"
	"testing"
	"time"

	"github.com/piyapong-mun/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomNumber(-1000, 1000),
	}

	entry, err := TestQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	assert.Equal(t, arg.AccountID, entry.AccountID)
	assert.Equal(t, arg.Amount, entry.Amount)
	assert.NotZero(t, entry.ID)
	assert.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	// Cleanup
	err := TestQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	clearAccount(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	entry2, err := TestQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	assert.Equal(t, entry1.ID, entry2.ID)
	assert.Equal(t, entry1.AccountID, entry2.AccountID)
	assert.Equal(t, entry1.Amount, entry2.Amount)
	assert.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	// Cleanup
	err = TestQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	clearAccount(t, account.ID)
}

func TestUpdateEntry(t *testing.T) {
	account1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, account1)

	account2 := createRandomAccount(t)

	arg := UpdateEntryParams{
		ID:        entry1.ID,
		AccountID: account2.ID,
		Amount:    util.RandomNumber(-1000, 1000),
	}

	err := TestQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)

	entry2, err := TestQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	assert.Equal(t, entry1.ID, entry2.ID)
	assert.Equal(t, arg.AccountID, entry2.AccountID)
	assert.Equal(t, arg.Amount, entry2.Amount)
	assert.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	// Cleanup
	err = TestQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	clearAccount(t, account1.ID)
	clearAccount(t, account2.ID)
}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	err := TestQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := TestQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)

	// Cleanup
	clearAccount(t, account.ID)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	var lastEntries []Entry
	for i := 0; i < 10; i++ {
		lastEntries = append(lastEntries, createRandomEntry(t, account))
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := TestQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

	// Cleanup
	for _, entry := range lastEntries {
		err = TestQueries.DeleteEntry(context.Background(), entry.ID)
		require.NoError(t, err)
	}
	clearAccount(t, account.ID)
}
