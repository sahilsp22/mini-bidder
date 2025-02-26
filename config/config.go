package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

const (
	CACHE_TIMEOUT = 60
	CACHE_UPDATE_INTERVAL = 30
)

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type Memcache struct {
	Host string
	Port string
}

type Creative struct {
	AdID 			string 	`json:"AdID"`
	Height  		string 	`json:"Height"`
	Width 			string 	`json:"Width"`
	AdType 			string 	`json:"AdType"`
	CreativeDetails string	`json:"CreativeDetails"`
}

func GetPGConfig() (*Postgres,error) {

	err := godotenv.Load("/Users/sahilpatil/mini-bidder/config/db.env")
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Postgres{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DB: os.Getenv("DB_NAME"),
	},nil
}

func GetMcCConfig() (*Memcache,error) {
	err := godotenv.Load("/Users/sahilpatil/mini-bidder/config/db.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Memcache{
		Host: os.Getenv("MC_HOST"),
		Port: os.Getenv("MC_PORT"),
	},nil
}

var query string = `CREATE_DB=create database test;
\c test
create user sp with password 1234;
grant all on schema public to sp;
create table t1(name varchar(26),age int);
insert into t1 values('sp',21),('ab',22);`
// create table Creative_Details(adid varchar(20),height int, width int,adtype int,crtv_details varchar(20));
// insert into Creative_Details values('adtest1',100,100,1,'addetails'),insert into Creative_Details values('adtest2',100,50,2,'addetails');