package sqlc_custom

import (
	"context"
	"database/sql"
	"fmt"

	sqlc "github.com/stuneak/simplebank/db/sqlc"
)

type Store interface {
	sqlc.Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
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

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    sqlc.Transfer `json:"transfer"`
	FromAccount sqlc.Account  `json:"from_account"`
	ToAccount   sqlc.Account  `json:"to_account"`
	FromEntry   sqlc.Entry    `json:"from_entry"`
	ToEntry     sqlc.Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *sqlc.Queries) error {

		var err error

		result.Transfer, err = q.CreateTransfer(ctx, sqlc.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// fix order to avoid deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}

		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *sqlc.Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 sqlc.Account, account2 sqlc.Account, err error) {
	account1, err = q.AddAccountBalance(ctx, sqlc.AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, sqlc.AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	if err != nil {
		return
	}

	return
}
