package internal

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	sqlc "github.com/stuneak/simplebank/db/sqlc"
	"github.com/stuneak/simplebank/util"
)

func createRandomAccount(t *testing.T) sqlc.Account {
	arg := sqlc.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)
	acc1, err := testQueries.GetAccount(context.Background(), acc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc1)

	require.Equal(t, acc.Owner, acc1.Owner)
	require.Equal(t, acc.ID, acc1.ID)
	require.Equal(t, acc.Currency, acc1.Currency)
	require.Equal(t, acc.Balance, acc1.Balance)
	require.WithinDuration(t, acc.CreatedAt, acc.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	arg := sqlc.UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, acc.Owner, acc1.Owner)
	require.Equal(t, acc.ID, acc1.ID)
	require.Equal(t, acc.Currency, acc1.Currency)
	require.Equal(t, arg.Balance, acc.Balance)
	require.WithinDuration(t, acc.CreatedAt, acc.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	// fmt.Print(acc2)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := sqlc.ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
