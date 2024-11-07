package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

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

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

type BuyItemTxParams struct {
	UserID int64 `json:"user_id"`
	ItemID int64 `json:"item_id"`
}

type BuyItemTxResult struct {
	BoughtItem  Item  `json:"bought_item"`
	OrderDetail Order `json:"order_detail"`
}

func (store *Store) BuyItemTx(ctx context.Context, arg BuyItemTxParams) (BuyItemTxResult, error) {
	var result BuyItemTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var item Item

		item, err = q.GetItem(ctx, arg.ItemID)
		if err != nil {
			log.Printf("Can't get item to buy: %v", err)
			return err
		}

		if item.Quantity < 1 {
			err = fmt.Errorf("Item id %d has %d item left", item.ID, item.Quantity)
			log.Printf("failed to buy: %v", err)
			return err
		}

		result.BoughtItem, err = q.UpdateItemQuantity(ctx, UpdateItemQuantityParams{
			Quantity: 1,
			ID:       arg.ItemID,
		})

		if err != nil {
			log.Printf("Error occured when updating item quant: %v", err)
			return err
		}

		result.OrderDetail, err = q.CreateOrder(ctx, CreateOrderParams{
			ItemID: arg.ItemID,
			UserID: arg.UserID,
		})

		if err != nil {
			log.Printf("Error occured when creating order: %v", err)
			return err
		}

		return nil
	})

	return result, err
}
