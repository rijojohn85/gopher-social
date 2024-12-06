package main

import (
	"github.com/rijojohn85/social/internal/db"
	"github.com/rijojohn85/social/internal/env"
	"github.com/rijojohn85/social/internal/store"
	"log"
)

func main() {

	addr := env.GetString(
		"DB_ADDR",
		"postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable",
	)
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	seedStore := store.NewStorage(conn)
	db.Seed(seedStore, conn)
}
