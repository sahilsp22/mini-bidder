package main
import (
	"fmt"
	// "log"
	"context"
	// "os"
	
	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
)

func main() {

	cntlog := logger.InitLogger(logger.CONTROLLER)

	cfg,err := config.GetPGConfig()
	if err!=nil {
		cntlog.Fatal(err)
	}
	fmt.Println(cfg)

	mcfg,err := config.GetMcCConfig()
	if err!=nil {
		cntlog.Fatal(err)
	}

	mc,err := db.NewMcClient(mcfg)
	if err!=nil {
		cntlog.Fatal(err)
	}
	// fmt.Println(mc)
	defer mc.Close()

	pg,err := db.NewClient(cfg)
	if err!=nil {
		cntlog.Fatal(err)
	}
	// fmt.Println(pg)
	rows,err := pg.Query(context.Background(), "SELECT * FROM Creative_Details;")
	if err!=nil {
		cntlog.Fatal(err)
	}
	var Creatives []config.Creative
	for rows.Next() {
		var crtv config.Creative
		err = rows.Scan(&crtv.AdID, &crtv.Height, &crtv.Width, &crtv.AdType, &crtv.CreativeDetails)
		if err != nil {
			cntlog.Fatal(err)
		}
		fmt.Println(crtv)
		Creatives = append(Creatives,crtv)
	}	
	rows.Close()

	for _,crtv := range Creatives {
		mc.Set(crtv.AdID,crtv)
	}

	if err = rows.Err(); err != nil {
        cntlog.Fatal(err)
    }

	crtv,err := mc.Get("adtest1")
	fmt.Println(crtv)

	defer pg.Close()


	return
}