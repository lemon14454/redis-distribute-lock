-- name: CreateItem :one
INSERT INTO items (
  name,
  quantity
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetItem :one
SELECT * FROM items
where id = $1 limit 1;

-- name: GetItemForUpdate :one
SELECT * FROM items
where id = $1 limit 1
FOR NO KEY UPDATE;

-- name: UpdateItemQuantity :one
UPDATE items
SET quantity = quantity - sqlc.arg(quantity)
WHERE id = sqlc.arg(id)
RETURNING *;
