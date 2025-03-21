package main
import (
	"fmt"
	// "context"
	// "time"
	// "log"
	// "os"
	
	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
	"github.com/sahilsp22/mini-bidder/utils"
)

func main() {

	cntlog := logger.InitLogger(logger.CONTROLLER)

	// Load configs
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

	pg,err := db.NewClient(cfg)
	if err!=nil {
		cntlog.Fatal(err)
	}

	// Get controller client
	controller,err := utils.NewController(pg,mc,cntlog)
	if err!=nil {
		cntlog.Fatal(err)
	}

	// Start controller
	controller.Start()

	defer pg.Close()
	defer mc.Close()

	return
}