package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const UserID = int64(69)

func TestBuyItem(t *testing.T) {
	store := NewStore(testDB)

	item := createRandomItem(t)
	fmt.Println("Quantity Before Buy:", item.Quantity)

	// run n concurrent tranfer transaction
	n := 5

	errs := make(chan error)
	results := make(chan BuyItemTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.BuyItemTx(ctx, BuyItemTxParams{
				UserID: UserID,
				ItemID: item.ID,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		resultItem := result.BoughtItem
		require.NotEmpty(t, item)
		require.Equal(t, item.ID, resultItem.ID)
		require.Equal(t, item.Quantity-(int64(i+1)), resultItem.Quantity)
		require.NotZero(t, resultItem.ID)

		order := result.OrderDetail
		_, err = store.GetOrder(context.Background(), order.ID)
		require.NoError(t, err)
		require.NotEmpty(t, order)
		require.Equal(t, order.ItemID, item.ID)
		require.Equal(t, order.UserID, UserID)
		require.NotZero(t, order.ID)
		require.NotZero(t, order.CreatedAt)
	}

	afterItem, err := store.GetItem(context.Background(), item.ID)
	fmt.Println("Quantity After Buy:", afterItem.Quantity)
	require.NoError(t, err)
	require.Equal(t, afterItem.ID, item.ID)
	require.Equal(t, afterItem.Quantity, item.Quantity-int64(n))

}
