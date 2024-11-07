package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"server/api"
	"server/cache"
	db "server/db/sqlc"
	"server/util"

	_ "github.com/lib/pq"
)

func main() {

	port := flag.String("port", "8080", "The port the server will listen on")
	flag.Parse()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannnot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	conn.SetMaxOpenConns(30)

	if err != nil {
		log.Fatal("failed to connect DB: ", err)
	}

	store := db.NewStore(conn)
	redis := cache.NewRedisClient()
	server := api.NewServer(store, redis)

	err = server.Start(fmt.Sprintf("127.0.0.1:%s", *port))

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
