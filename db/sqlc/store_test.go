package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrasnferTx(t *testing.T) {
	store := NewStore(testDB)

	payerAccount := createRandomAccount(t)
	collectorAccount := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: payerAccount.ID,
				ToAccountID:   collectorAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, payerAccount.ID, transfer.FromAccountID)
		require.Equal(t, collectorAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, payerAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, collectorAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, payerAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, collectorAccount.ID, toAccount.ID)

		diff1 := payerAccount.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - collectorAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), payerAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), collectorAccount.ID)
	require.NoError(t, err)

	require.Equal(t, payerAccount.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, collectorAccount.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTrasnferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	payerAccount := createRandomAccount(t)
	collectorAccount := createRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		payerAccountID := payerAccount.ID
		collectorAccountID := collectorAccount.ID

		if i%2 == 1 {
			payerAccountID = collectorAccount.ID
			collectorAccountID = payerAccount.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: payerAccountID,
				ToAccountID:   collectorAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := store.GetAccount(context.Background(), payerAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), collectorAccount.ID)
	require.NoError(t, err)

	require.Equal(t, payerAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, collectorAccount.Balance, updatedAccount2.Balance)

}
