package db

import (
	"fmt"
	// "log"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sahilsp22/mini-bidder/db/config"
)

type PgClient struct{
	cl *pgxpool.Pool
}

func NewClient(cfg *config.Config) (*PgClient, error) {
	conf,err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",cfg.User,cfg.Password,cfg.Host,cfg.Port,cfg.DB))

	if err!=nil {
		// log.Fatal(err)
		return nil,err
	}
	
	cl, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		// log.Fatal(err)
		return nil,err
	}

	fmt.Println("Connected to DB")

	return &PgClient{cl:cl}, nil
}


func (pg *PgClient) Query(ctx context.Context, query string,args ...interface{}) (pgx.Rows, error) {
	rows,err := pg.cl.Query(ctx, query,args...)
	if err!=nil {
		// log.Fatal(err)
		return nil,err
	}
	return rows,nil
}

func (pg *PgClient) Close() {
	pg.cl.Close()
}

func (pg *PgClient) Exec(ctx context.Context, query string,args ...interface{}) (error) {
	_,err := pg.cl.Exec(ctx, query,args...)
	if err!=nil {
		// log.Fatal(err)
		return err
	}
	return nil
}
