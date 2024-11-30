package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-backend-projects/blogging_platform_api/cmd/api"
	"github.com/go-backend-projects/blogging_platform_api/configs"
	"github.com/go-backend-projects/blogging_platform_api/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	cfg := configs.Envs.DBConfig

	db, err := db.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initDB(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initDB(db *pgxpool.Pool) {

	err := db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database %v", err)
	}

	log.Println("Successful connect to DB")
}
