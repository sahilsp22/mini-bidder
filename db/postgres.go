package db

import (
	"fmt"
	// "log"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
	
)

type PgClient struct{
	cl *pgxpool.Pool
}

var PgInstance *PgClient

var pglog *logger.Logger
func init() {
    pglog = logger.InitLogger(logger.POSTGRES)

}

// Returns a Postgres client
func NewClient(cfg *config.Postgres) (*PgClient, error) {

	conf,err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",cfg.User,cfg.Password,cfg.Host,cfg.Port,cfg.DB))
	if err!=nil {
		pglog.Fatal(err)
		return nil,err
	}
	
	cl, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		pglog.Fatal(err)
		return nil,err
	}

	PgInstance = &PgClient{cl:cl}
	pglog.Print("DB Initialized")
	pglog.Print("Connected to DB")
	
	return PgInstance,nil
}

func GetPGInstance() *PgClient {
	return PgInstance
}

// Performs a query on PG 
func (pg *PgClient) Query(ctx context.Context, query string,args ...interface{}) (pgx.Rows, error) {
	rows,err := pg.cl.Query(ctx, query,args...)
	if err!=nil {
		// pglog.Print(err)
		return nil,fmt.Errorf("error querying: %v", err)
	}
	return rows,nil
}

// Query without a return values
func (pg *PgClient) Exec(ctx context.Context, query string,args ...interface{}) (error) {
	_,err := pg.cl.Exec(ctx, query,args...)
	if err!=nil {
		// pglog.Print(err)
		return fmt.Errorf("error executing query: %v", err)
	}
	return nil
}

func (pg *PgClient) Close() {
	pg.cl.Close()
}