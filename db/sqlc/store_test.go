package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1, _ := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "MiuMiu",
		Balance:  1000000,
		Currency: "COP",
	})

	account2, _ := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "JB",
		Balance:  1000,
		Currency: "EUR",
	})

	concurrentTransactionsTotal := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < concurrentTransactionsTotal; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < concurrentTransactionsTotal; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer.
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// Check entries.
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// Check accounts.
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		// Check Balance.
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := account2.Balance - result.ToAccount.Balance
		require.Equal(t, diff1, -diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, ..., n * amount.

		k := int(diff1 / amount)
		require.True(t, k > 0 && k <= concurrentTransactionsTotal)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated Balances.
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(concurrentTransactionsTotal)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(concurrentTransactionsTotal)*amount, updateAccount2.Balance)
}
