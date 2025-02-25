package main
import (
	"fmt"
	"log"
	
	"minibidder/db"
)

func main() {
	pg,err := db.NewClient()
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(pg)
	rows,err := pg.Query(context.Background(), "SELECT * FROM t1")
	if err!=nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}	
	defer pg.Close()
	return
}