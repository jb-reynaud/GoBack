package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueries_CreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "MiuMiu",
		Balance:  1000000,
		Currency: "COP",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestQueries_GetAccount(t *testing.T) {
	account1, _ := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "MiuMiu",
		Balance:  1000000,
		Currency: "COP",
	})
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestQueries_ListAccounts(t *testing.T) {
	_, _ = testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "MiuMiu",
		Balance:  1000000,
		Currency: "COP",
	})
	_, _ = testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "JB",
		Balance:  1000,
		Currency: "EUR",
	})
	_, _ = testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "Katha",
		Balance:  100,
		Currency: "CAD",
	})

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  2,
		Offset: 1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 2)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

func TestQueries_UpdateAccount(t *testing.T) {
	account1, _ := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "MiuMiu",
		Balance:  1000000,
		Currency: "COP",
	})

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: 2000000,
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestQueries_DeleteAccount(t *testing.T) {
	account1, _ := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "MiuMiu",
		Balance:  1000000,
		Currency: "COP",
	})

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}
