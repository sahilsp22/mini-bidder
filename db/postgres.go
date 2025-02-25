package db

import (
	"fmt"
	"log"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgClient *pgxpool.Pool

func NewClient() (*PgClient, error) {
	conf,err := pgxpool.ParseConfig("postgres://sp:1234@192.168.64.2:30336/test")

	if err!=nil {
		// log.Fatal(err)
		return nil,err
	}
	
	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		// log.Fatal(err)
		return nil,err
	}

	fmt.Println("Connected to DB")

	return &pool, nil
}

