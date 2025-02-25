package main
import (
	"fmt"
	"log"
	"context"
	
	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/db/config"
)

func main() {

	cfg,err := config.GetPGConfig()
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	pg,err := db.NewClient(&cfg)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(pg)
	rows,err := pg.Query(context.Background(), "SELECT * FROM t1")
	if err!=nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var yr string
		var name string
		err = rows.Scan(&name, &yr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(yr, name)
	}	

	if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
	defer rows.Close()
	defer pg.Close()
	return
}