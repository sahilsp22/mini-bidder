package main
import (
	"fmt"
	"context"
	"time"
	"sync"
	// "log"
	// "os"
	
	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
)

func main() {

	cntlog := logger.InitLogger(logger.CONTROLLER)

	wg := sync.WaitGroup{}
	defer wg.Wait()

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

	updateTicker := time.NewTicker(time.Second * config.CACHE_UPDATE_INTERVAL)

	wg.Add(1)
	go func() {
		for range updateTicker.C {
			start:=time.Now()
			updateCreatives()
			cntlog.Printf("Updated Creatives in %v",time.Since(start).Milliseconds())
		}
		defer wg.Done()
		defer updateTicker.Stop()
	}()

	go func() {
		for range updateTicker.C {
			start:=time.Now()
			updateAdvertisers()
			cntlog.Printf("Updated Advertiser Budgets in %v",time.Since(start).Milliseconds())
		}
		defer wg.Done()
		defer updateTicker.Stop()
	}()

	wg.Wait()
	defer pg.Close()

	return
}

func updateCreatives() {
	rows,err := pg.Query(context.Background(), "SELECT * FROM Creative_Details;")
	if err!=nil {
		logger.GetLogger(logger.CONTROLLER).Fatal(err)
	}
	var Creatives []config.Creative
	for rows.Next() {
		var crtv config.Creative
		err = rows.Scan(&crtv.AdID, &crtv.Height, &crtv.Width, &crtv.AdType, &crtv.CreativeDetails)
		if err != nil {
			logger.GetLogger(logger.CONTROLLER).Fatal(err)
		}
		// fmt.Println(crtv)
		Creatives = append(Creatives,crtv)
	}	
	if err = rows.Err(); err != nil {
		logger.GetLogger(logger.CONTROLLER).Fatal(err)
	}
	rows.Close()
	
	for _,crtv := range Creatives {
		err:=mc.Set(crtv.AdID,crtv)
		if err!=nil{
			cntlog.Print(err)
		}
	}
	
	// var crtv config.Creative
	// err = mc.Get("adtest3",&crtv)
	// if err!=nil{
	// 	logger.GetLogger(logger.CONTROLLER).Print(err)
	// }
	// fmt.Println(crtv)
}

func updateAdvertisers() {
	
}