package main

import (
	// "fmt"
	// "time"
	// "context"
	// "log"
	// "os"

	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
	// "github.com/sahilsp22/mini-bidder/utils"
	"github.com/sahilsp22/mini-bidder/server"
	"github.com/sahilsp22/mini-bidder/bid"
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

	// cntrlClient,err := utils.NewController(pg,mc)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	m := bid.NewMatcher(pg,mc,logger)

	srvr := server.Server{}
	metricsrv := server.Server{}

	routes := []server.Route{
		{
			Path: "/bid",
			Handler: m.Handler,
		},
	}

	metricroutes := []server.Route{
		{
			Path: "/metrics",
			Handler: m.Metrics().ServeHTTP,
		},
	}

	srvr.AddRoutes(routes)
	metricsrv.AddRoutes(metricroutes)

	go func(){
		metricsrv.Listen(config.METRICS_SERVER_PORT)
		// logger.Print("Metrics server listening on port : ",8080)
		}()
		
	srvr.Listen(config.BIDDER_SERVER_PORT)
	// logger.Print("Bid server listening on port : ",3333)

	defer mc.Close()
	defer pg.Close()

	return
}

// var bd utils.Budget
// err = mc.Get("advtest1",&bd)
// if err!=nil {
// 	logger.Print(err)
// } else {
// 	fmt.Println(bd)
// }
	
// cntrl,err := utils.NewController(pg,mc)
// err=cntrl.UpdateAdvBudget("advtest1")
// if err!= nil {
// 	logger.Print(err)
// }
// time.Sleep(10*time.Second)

// var b utils.Budget
// err = mc.Get("advtest1",&b)
// if err!=nil {
// 	logger.Print(err)
// } else {
// 	fmt.Println(b)
// }
