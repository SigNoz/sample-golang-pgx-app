package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

func connectDB() {
	var err error
	connConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	fmt.Println("connConfig.ConnConfig.Tracer: ", connConfig.ConnConfig.Tracer)
	conn, err = pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}
