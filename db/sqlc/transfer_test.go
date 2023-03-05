package db

import (
	"context"
	"testing"
	"time"

	util "github.com/josephsae/bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, collectorAccount, PayerAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: PayerAccount.ID,
		ToAccountID:   collectorAccount.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	collectorAccount := createRandomAccount(t)
	PayerAccount := createRandomAccount(t)
	createRandomTransfer(t, collectorAccount, PayerAccount)
}

func TestGetTransfer(t *testing.T) {
	collectorAccount := createRandomAccount(t)
	PayerAccount := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, collectorAccount, PayerAccount)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	collectorAccount := createRandomAccount(t)
	PayerAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, collectorAccount, PayerAccount)
	}

	arg := ListTransfersParams{
		FromAccountID: PayerAccount.ID,
		ToAccountID:   collectorAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
