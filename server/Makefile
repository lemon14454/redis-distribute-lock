postgres:
	docker run --name scale_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it scale_db createdb --username=root --owner=root scale

dropdb:
	docker exec -it scale_db dropdb --username=root --owner=root scale

cleardb:
	docker exec -it scale_db psql -U root -d scale -c "TRUNCATE TABLE items, orders RESTART IDENTITY;"

item:
	docker exec -it scale_db psql -U root -d scale -c "INSERT INTO items(name, quantity) VALUES ('TEST', 30);"

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/scale?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/scale?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

GO_RUN = go run main.go
PORT1 = 8080
PORT2 = 8081
PORT3 = 8082

server:
	go run main.go

server1:
	$(GO_RUN) -port=$(PORT1)

server2:
	$(GO_RUN) -port=$(PORT2)

server3:
	$(GO_RUN) -port=$(PORT3)

.PHONY: postgres createdb dropdb cleardb migrateup migratedown sqlc test server server1 server2 server3

