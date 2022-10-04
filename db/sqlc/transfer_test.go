package db

import (
	"context"
	"testing"
	"time"

	"github.com/harshva/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1 Account, account2 Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	transfer1 := createRandomTransfer(t, acc1, acc2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

}

func TestListGetTransfers(t *testing.T) {
	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, acc1, acc2)
		createRandomTransfer(t, acc2, acc1)
	}

	arg := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == acc1.ID || transfer.ToAccountID == acc1.ID)
	}
}
