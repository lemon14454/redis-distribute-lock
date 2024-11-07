package db

import (
	"context"
	"server/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomItem(t *testing.T) Item {
	arg := CreateItemParams{
		Name:     util.RandomItem(),
		Quantity: util.RandomQuantity(),
	}

	item, err := testQueries.CreateItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	require.Equal(t, arg.Name, item.Name)
	require.Equal(t, arg.Quantity, item.Quantity)

	require.NotZero(t, item.ID)

	return item
}

func TestCreateItem(t *testing.T) {
	createRandomItem(t)
}

func TestGetItem(t *testing.T) {

	item1 := createRandomItem(t)
	item2, err := testQueries.GetItem(context.Background(), item1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item1.ID, item2.ID)
	require.Equal(t, item1.Name, item2.Name)
	require.Equal(t, item1.Quantity, item2.Quantity)
}
