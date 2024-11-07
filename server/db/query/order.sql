-- name: CreateOrder :one
INSERT INTO orders (
  item_id,
  user_id
) VALUES (
  $1, $2
) RETURNING *;


-- name: GetOrder :one
SELECT * FROM orders
where id = $1 limit 1;
