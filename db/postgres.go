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

func NewClient(cfg *config.Postgres) (*PgClient, error) {
	// InitDB()
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
	// PgInstance.InitDB()
	
	return PgInstance,nil
}

func GetPGInstance() *PgClient {
	return PgInstance
}

func InitDB(){
	
	conf,err := pgxpool.ParseConfig(fmt.Sprintf("postgres://postgres:password@192.168.64.2:30336/postgres"))
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return 
	}
	cl, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return 
	}

	_,err = cl.Exec(context.Background(), `DROP DATABASE IF EXISTS test;`)
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	_,err = cl.Exec(context.Background(), `CREATE DATABASE test;`)
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	_,err = cl.Exec(context.Background(), `CREATE ROLE sp SUPERUSER LOGIN PASSWORD '1234';`)
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	_,err = cl.Exec(context.Background(), `GRANT ALL ON SCHEMA public TO sp;`)
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	_,err = cl.Exec(context.Background(), `GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO sp;`)
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}

	tcl,err:=pgx.Connect(context.Background(), "postgres://sp:1234@192.168.64.2:30336/test")
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	_,err = tcl.Exec(context.Background(), `CREATE TABLE Creative_Details(adid varchar(20),height int, width int,adtype int,crtv_details varchar(20));`)
	if err!=nil {
		fmt.Println(err)
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	_,err = tcl.Exec(context.Background(), `INSERT INTO Creative_Details values('adtest1',100,100,1,'addetails'),('adtest2',100,50,2,'addetails');`)
	if err!=nil {
		pglog.Fatalf(fmt.Sprintf("Error Initialising DB: %s",err))
		return
	}
	defer cl.Close()
	defer tcl.Close(context.Background())
}


func (pg *PgClient) Query(ctx context.Context, query string,args ...interface{}) (pgx.Rows, error) {
	rows,err := pg.cl.Query(ctx, query,args...)
	if err!=nil {
		// pglog.Print(err)
		return nil,fmt.Errorf("error querying: %v", err)
	}
	return rows,nil
}

func (pg *PgClient) Close() {
	pg.cl.Close()
}

func (pg *PgClient) Exec(ctx context.Context, query string,args ...interface{}) (error) {
	_,err := pg.cl.Exec(ctx, query,args...)
	if err!=nil {
		// pglog.Print(err)
		return fmt.Errorf("error executing query: %v", err)
	}
	return nil
}
