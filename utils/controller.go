package utils

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

type Creative struct {
	AdID 			string 	`json:"AdID"`
	Height  		string 	`json:"Height"`
	Width 			string 	`json:"Width"`
	AdType 			string 	`json:"AdType"`
	CreativeDetails string	`json:"CreativeDetails"`
	AdvertiserID    string
}

type Budget struct {
	AdvID 		string 	`json:"AdvID"`
	Budget 		string 	`json:"totalBudget"`
	CPM 		string 	`json:"cpm"`
	RemBudget 	string	`json:"remBudget"`
}

type Controller struct {
	pg *db.PgClient
	mc *db.MCacheClient
	logger *logger.Logger
}


var (
	// controller *Controller
	crtvLock *sync.Mutex
	bdgtLock *sync.RWMutex
)

func init() {
	crtvLock = &sync.Mutex{}
	bdgtLock = &sync.RWMutex{}
}

func NewController(pg *db.PgClient, mc *db.MCacheClient, logger *logger.Logger) (*Controller,error) {

	controller := &Controller{
		pg:pg,
		mc:mc,
		logger:logger,
	}

	return controller,nil
}

// func GetControllerInstance() *Controller {
// 	fmt.Println(controller)
// 	return controller
// }

func (c *Controller) Start() {

	// logger:=logger.InitLogger(logger.CONTROLLER)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	updateCrtvTicker := time.NewTicker(time.Second * config.CACHE_UPDATE_INTERVAL)
	updateBudgetTicker := time.NewTicker(time.Second * config.CACHE_UPDATE_INTERVAL)

	wg.Add(1)
	go func() {
		for range updateCrtvTicker.C {
			start:=time.Now()
			err:=c.updateCreatives()
			if err!=nil {
				c.logger.Print("Error updating Creatives: ",err)
				continue
			}
			c.logger.Print("Updated Creatives in: ",time.Since(start).Milliseconds())
		}
		wg.Done()
		updateCrtvTicker.Stop()
	}()

	wg.Add(1)
	go func() {
		for range updateBudgetTicker.C {
			start:=time.Now()
			err:=c.updateAdvertisers()
			if err!=nil {
				c.logger.Print("Error updating Budgets: ",err)
				continue
			}
			c.logger.Print("Updated Budgets in: ",time.Since(start).Milliseconds())
		}
		wg.Done()
		updateBudgetTicker.Stop()
	}()

	wg.Wait()

	return
}

func (c *Controller) updateCreatives() error{
	rows,err := c.pg.Query(context.Background(), ALL_CREATIVES_QUERY)
	if err!=nil {
		return fmt.Errorf("Error reading creatives: %v",err)
	}
	var Creatives []*Creative
	for rows.Next() {
		var crtv Creative
		err = rows.Scan(&crtv.AdID, &crtv.Height, &crtv.Width, &crtv.AdType, &crtv.CreativeDetails, &crtv.AdvertiserID)
		if err != nil {
			return fmt.Errorf("Error scanning Creative rows: %v",err)
		}
		// fmt.Println(crtv)
		Creatives = append(Creatives,&crtv)
	}	
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Error scanning Creative rows: %v",err)
	}
	rows.Close()
	
	crtvLock.Lock()
	for _,crtv := range Creatives {
		err:=c.mc.Set(crtv.AdID,crtv)
		if err!=nil{
			return err
		}
	}
	crtvLock.Unlock()
	return nil
	// var crtv config.Creative
	// err = mc.Get("adtest3",&crtv)
	// if err!=nil{
	// 	logger.GetLogger(logger.CONTROLLER).Print(err)
	// }
	// fmt.Println(crtv)
}

func (c *Controller) updateAdvertisers() error {

	bdgtLock.RLock()
	rows,err := c.pg.Query(context.Background(), ALL_BUDGET_QUERY)
	bdgtLock.RUnlock()

	if err!=nil {
		return fmt.Errorf("Error scanning Budget rows: %v",err)
	}

	var Budgets []*Budget
	for rows.Next() {
		var bdgt Budget
		err = rows.Scan(&bdgt.AdvID, &bdgt.Budget, &bdgt.CPM, &bdgt.RemBudget)
		if err != nil {
			return fmt.Errorf("Error scanning Budget rows: %v",err)
		}
		// fmt.Println(bdgt)
		Budgets = append(Budgets,&bdgt)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("Error scanning Budget rows: %v",err)
	}
	rows.Close()
	
	bdgtLock.Lock()
	defer bdgtLock.Unlock()
	for _,budget := range Budgets {
		err:=c.mc.Set(budget.AdvID,budget)
		if err!=nil{
			return err
		}
	}

	// var b Budget
	// err = c.mc.Get("advtest1",&b)
	// if err!=nil{
	// 	return fmt.Errorf("Error getting Budget: %v",err)
	// }
	// fmt.Println(b)
	return nil
}

func (c *Controller) UpdateAdvBudget(AdvID string) (error) {
	// fmt.Println("updating....")
	bdgtLock.Lock()
	defer bdgtLock.Unlock()
    query := UPDATE_ADV_BUDGET_QUERY
    err := c.pg.Exec(context.Background(), query, AdvID)
	if err!=nil {
		return fmt.Errorf("Error updating Budget: AdvID: %v : %v",AdvID,err)
	}
	
	go func() error {
		err:=c.updateAdvertisers()
		if err!=nil {
			return err
		}
		// fmt.Println("Updated Budget for Adv")
		return nil
	}()
	return nil
}

func CreativeLock() {
	crtvLock.Lock()
}
func CreativeUnLock() {
	crtvLock.Unlock()
}