package db

import (
	"context"
	"database/sql"

	// import uuid
	"github.com/google/uuid"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type StoreProduction struct {
	*Queries
	DB *sql.DB
}

func NewStore(db *sql.DB) *StoreProduction {
	return &StoreProduction{
		Queries: New(db),
		DB:      db,
	}
}

func (store *StoreProduction) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *StoreProduction) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err error

	err = store.execTx(ctx, func(q *Queries) error {
		// Create transfer
		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create Entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: Update account's balance

		// Prevent Deadlock
		// With this method, even if two transactions try to update the same account's balance at the same time,
		// only one transaction will be able to update the account's balance at a time.

		// If we dont follow the same order to update account's balance
		// TX1 update Account1 then [Account1 is Locked]
		// TX2 update Account2 then [Account2 is Locked]
		// TX1 update Account2 then [Account2 is Locked by TX2] -> TX1 cannot update Account1
		// TX2 update Account1 then [Account1 is Locked by TX1] -> TX2 cannot update Account2 -> Deadlock

		// If we follow the same order to update account's balance
		// TX1 update Account1 then [Account1 is Locked]
		// TX2 update Account1 then [Account1 is Locked by TX1] -> wait for TX1 to finish its transaction
		// TX1 update Account2 then [Account2 is Locked] -> Now Account did not lock by TX2
		// TX2 update Account2 then [Account2 is Unlocked by TX1] -> No Deadlock

		// The key is to always update account's balance in the same order
		if uuidHashToInt(arg.FromAccountID) <= uuidHashToInt(arg.ToAccountID) { // Account A หรือ Account B มากกว่า
			// ถ้า Account A น้อยกว่า Account B ให้ update Account A ก่อน แล้ว update Account B
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			// ถ้า Account A น้อยกว่า Account B ให้ update Account A ก่อน แล้ว update Account B
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func addMoney(ctx context.Context, q *Queries, account1ID uuid.UUID, amount1 int64, account2ID uuid.UUID, amount2 int64) (Account, Account, error) {

	// Two way of implement Update Account Balance in Postgres
	// [First Method] use `SELECT ... FOR UPDATE` and then `UPDATE`
	// [Second Method] use `UPDATE ... SET balance = balance + ?`
	// But [Second Method] is more efficient because it doesn't need to read old balance

	// [First Method] use `SELECT ... FOR UPDATE` and then `UPDATE`
	// 1. Lock the row manually (Confirm that no changes will be made to the row)
	// 2. then read old balance to use it in Update (using old balance - amount)
	// The lock will last until the transaction is committed

	// Block other transaction from updating the account's balance until this transaction is committed
	acc1, err := q.GetAccountForUpdate(ctx, account1ID)
	if err != nil {
		return Account{}, Account{}, err
	}
	acc1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:      acc1.ID,
		Balance: acc1.Balance + amount1,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	// [Second Method] use `UPDATE ... SET balance = balance + ?`
	// Let database lock the row automatically, when we use UPDATE database will lock the row automatically
	// In this method we don't need to read old balance, because we update `balance + amount` directly
	// The lock will last until the transaction is committed
	acc2, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	if err != nil {
		return Account{}, Account{}, err
	}
	return acc1, acc2, nil
}

func uuidHashToInt(uuid uuid.UUID) int64 {
	// sum all ASCII Code of uuid
	var sum int64
	for _, b := range uuid {
		sum += int64(b)
	}
	return sum
}
