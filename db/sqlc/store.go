package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
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

func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from an account to another.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = addAmount(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)

		return nil
	})

	return result, err
}

func addAmount(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (fromAccount Account, toAccount Account, err error) {
	if fromAccountID < toAccountID {
		// Prevent deadlock.
		toAccount, fromAccount, err = addAmount(ctx, q, toAccountID, fromAccountID, -amount)

		return fromAccount, toAccount, err
	}

	fromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     fromAccountID,
		Amount: -amount,
	})
	if err != nil {
		return
	}

	toAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     toAccountID,
		Amount: amount,
	})
	if err != nil {
		return
	}

	return fromAccount, toAccount, err
}
