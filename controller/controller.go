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

	pg,err := db.NewClient(cfg)
	if err!=nil {
		cntlog.Fatal(err)
	}
	// fmt.Println(pg)

	updateCrtvTicker := time.NewTicker(time.Second * config.CACHE_UPDATE_INTERVAL)
	updateBudgetTicker := time.NewTicker(time.Second * config.CACHE_UPDATE_INTERVAL)

	wg.Add(1)
	go func() {
		for range updateCrtvTicker.C {
			start:=time.Now()
			err:=updateCreatives()
			if err!=nil {
				cntlog.Print("Error updating Creatives: ",err)
				continue
			}
			cntlog.Print("Updated Creatives in %v",time.Since(start).Milliseconds())
		}
		defer wg.Done()
		defer updateCrtvTicker.Stop()
	}()

	wg.Add(1)
	go func() {
		for range updateBudgetTicker.C {
			start:=time.Now()
			err:=updateAdvertisers()
			if err!=nil {
				cntlog.Print("Error updating Adv Budgets: ",err)
			}
			cntlog.Print("Updated Advertiser Budgets in %v",time.Since(start).Milliseconds())
		}
		defer wg.Done()
		defer updateBudgetTicker.Stop()
	}()

	wg.Wait()
	defer pg.Close()
	defer mc.Close()

	return
}

func updateCreatives() error{
	rows,err := pg.Query(context.Background(), "SELECT * FROM Creative_Details;")
	if err!=nil {
		return fmt.Errorf("Error reading creatives: %v",err)
	}
	var Creatives []config.Creative
	for rows.Next() {
		var crtv config.Creative
		err = rows.Scan(&crtv.AdID, &crtv.Height, &crtv.Width, &crtv.AdType, &crtv.CreativeDetails)
		if err != nil {
			return fmt.Errorf("Error scanning Creative rows: %v",err)
		}
		// fmt.Println(crtv)
		Creatives = append(Creatives,crtv)
	}	
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Error scanning Creative rows: %v",err)
	}
	rows.Close()
	
	for _,crtv := range Creatives {
		err:=mc.Set(crtv.AdID,crtv)
		if err!=nil{
			return err
		}
	}
	return nil
	// var crtv config.Creative
	// err = mc.Get("adtest3",&crtv)
	// if err!=nil{
	// 	logger.GetLogger(logger.CONTROLLER).Print(err)
	// }
	// fmt.Println(crtv)
}

func updateAdvertisers() error {
	rows,err := pg.Query(context.Background(), "SELECT * FROM Budget;")
	if err!=nil {
		return fmt.Errorf("Error scanning Budget rows: %v",err)
	}
	var Budgets []config.Budget
	for rows.Next() {
		var bdgt config.Budget
		err = rows.Scan(&bdgt.AdvID, &bdgt.Budget, &bdgt.CPM, &bdgt.RemBudget)
		if err != nil {
			return fmt.Errorf("Error scanning Budget rows: %v",err)
		}
		// fmt.Println(bdgt)
		Budgets = append(Budgets,bdgt)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Error scanning Budget rows: %v",err)
	}
	rows.Close()
	
	for _,budget := range Budgets {
		err:=mc.Set(budget.AdvID,budget)
		if err!=nil{
			return err
		}
	}
	return nil
}