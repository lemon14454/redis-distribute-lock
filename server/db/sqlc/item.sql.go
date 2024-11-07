// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: item.sql

package db

import (
	"context"
)

const createItem = `-- name: CreateItem :one
INSERT INTO items (
  name,
  quantity
) VALUES (
  $1, $2
)
RETURNING id, name, quantity
`

type CreateItemParams struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, createItem, arg.Name, arg.Quantity)
	var i Item
	err := row.Scan(&i.ID, &i.Name, &i.Quantity)
	return i, err
}

const getItem = `-- name: GetItem :one
SELECT id, name, quantity FROM items
where id = $1 limit 1
`

func (q *Queries) GetItem(ctx context.Context, id int64) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItem, id)
	var i Item
	err := row.Scan(&i.ID, &i.Name, &i.Quantity)
	return i, err
}

const getItemForUpdate = `-- name: GetItemForUpdate :one
SELECT id, name, quantity FROM items
where id = $1 limit 1
FOR NO KEY UPDATE
`

func (q *Queries) GetItemForUpdate(ctx context.Context, id int64) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItemForUpdate, id)
	var i Item
	err := row.Scan(&i.ID, &i.Name, &i.Quantity)
	return i, err
}

const updateItemQuantity = `-- name: UpdateItemQuantity :one
UPDATE items
SET quantity = quantity - $1
WHERE id = $2
RETURNING id, name, quantity
`

type UpdateItemQuantityParams struct {
	Quantity int64 `json:"quantity"`
	ID       int64 `json:"id"`
}

func (q *Queries) UpdateItemQuantity(ctx context.Context, arg UpdateItemQuantityParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, updateItemQuantity, arg.Quantity, arg.ID)
	var i Item
	err := row.Scan(&i.ID, &i.Name, &i.Quantity)
	return i, err
}