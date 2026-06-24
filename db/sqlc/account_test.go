package db

import (
	"context"

	"testing"

	"github.com/google/uuid"
	"github.com/piyapong-mun/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomNumber(1, 1000),
		Currency: util.RandomCurrency(),
	}

	account, err := TestQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)     // Check that there is no error
	assert.NotEmpty(t, account) // Check that the returned account is not empty

	assert.Equal(t, arg.Owner, account.Owner)       // Check that the owner matches
	assert.Equal(t, arg.Balance, account.Balance)   // Check that the balance matches
	assert.Equal(t, arg.Currency, account.Currency) // Check that the currency matches
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomNumber(1, 1000),
		Currency: util.RandomCurrency(),
	}

	account, err := TestQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	return account
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := TestQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := TestQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func clearAccount(t *testing.T, id uuid.UUID) {
	err := TestQueries.DeleteAccount(context.Background(), id)
	require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := TestQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)

	clearAccount(t, account1.ID)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := TestQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func selectAccount(t *testing.T, id uuid.UUID) Account {
	account, err := TestQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	return account
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountBalanceParams{
		ID:      account1.ID,
		Balance: util.RandomNumber(1, 1000),
	}

	_, err := TestQueries.UpdateAccountBalance(context.Background(), arg)

	account1_updated := selectAccount(t, account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account1_updated)

	assert.Equal(t, account1.ID, account1_updated.ID)
	assert.Equal(t, account1.Owner, account1_updated.Owner)
	assert.Equal(t, arg.Balance, account1_updated.Balance)
	assert.Equal(t, account1.Currency, account1_updated.Currency)
}
