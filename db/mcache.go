package db

import (
	"fmt"
	"encoding/json"
	// "log"
	// "context"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/sahilsp22/mini-bidder/config"
	"github.com/sahilsp22/mini-bidder/logger"
)


type MCacheClient struct{
	cl *memcache.Client
}

var McInstance *MCacheClient

var mclog *logger.Logger
func init() {
	mclog = logger.InitLogger(logger.MEMCACHE)
}

func NewMcClient(cfg *config.Memcache) (*MCacheClient, error) {
	mc := memcache.New(fmt.Sprintf("%s:%s",cfg.Host,cfg.Port))
	mclog.Print("Connected to Memcache")
	McInstance = &MCacheClient{cl:mc}
	return McInstance,nil
}

func (mc *MCacheClient) Set(key string, value interface{}) error {

	bs,err:=json.Marshal(value)
	if err!=nil {
		return fmt.Errorf("error marshalling value: %v", err)
	}
	// fmt.Println(string(bs))
	err =  mc.cl.Set(
			&memcache.Item{
			Key: key, 
			Value: bs,
			Expiration: config.CACHE_TIMEOUT,
		})
	if err!=nil {
		return fmt.Errorf("error setting key: %v", err)
	}
	return nil
}

func (mc *MCacheClient) Get(key string) (interface{},error) {
	item, err := mc.cl.Get(key)
	if err != nil {
		mclog.Print(err)
		if err == memcache.ErrCacheMiss {
			return nil,fmt.Errorf("key not found: %v", err)
		}
		if err == memcache.ErrMalformedKey {
			return nil,fmt.Errorf("malformed key: %v", err)
		}
		return nil,fmt.Errorf("error getting key: %v", err)
	}
	fmt.Println(string(item.Value))
	var crtv config.Creative
	err = json.Unmarshal(item.Value, &crtv)
	if err!=nil {
		return nil,fmt.Errorf("error unmarshalling value: %v", err)
	}
	return crtv,nil
}

func (mc *MCacheClient) Close() error{
	err := mc.cl.Close()
	if err!=nil {
		return fmt.Errorf("Client closed with error: %v", err)
	}
	return nil
}
func GetMcInstance() *MCacheClient {
	return McInstance
}