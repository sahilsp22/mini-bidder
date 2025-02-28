package main

import (
	"fmt"
	"time"
	// "context"
	// "log"
	// "os"

	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
	"github.com/sahilsp22/mini-bidder/utils"
)

func main() {

	logger := logger.InitLogger(logger.BIDDER)

	mcfg,err := config.GetMcCConfig()
	if err!=nil {
		logger.Fatal(err)
	}

	mc,err:= db.NewMcClient(mcfg)
	if err!=nil {
		logger.Fatal(err)
	}
	// mc:=db.GetMcInstance()

	// pg := db.GetPGInstance()
	pgcfg,err:=config.GetPGConfig()
	if err != nil {
		logger.Fatal(err)
	}

	pg, err := db.NewClient(pgcfg)
	if err != nil {
		logger.Fatal(err)
	}

	// for  {
		// fmt.Println(mc)
		// mc.Lock()
		var bd utils.Budget
		err = mc.Get("advtest1",&bd)
		if err!=nil {
			logger.Print(err)
		} else {
			fmt.Println(bd)
		}
		// mc.Unlock()
		// }
		
	cntrl,err := utils.NewController(pg,mc)
	err=cntrl.UpdateAdvBudget("advtest1")
	if err!= nil {
		logger.Print(err)
	}
	time.Sleep(10*time.Second)

	var b utils.Budget
	err = mc.Get("advtest1",&b)
	if err!=nil {
		logger.Print(err)
	} else {
		fmt.Println(b)
	}

	defer mc.Close()
	defer pg.Close()

	return
}
