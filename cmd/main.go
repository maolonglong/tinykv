package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/maolonglong/tinykv"
	"github.com/maolonglong/tinykv/config"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *tinykv.Group {
	return tinykv.NewGroup("scores", 2<<10, tinykv.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, peers []string, kv *tinykv.Group) {
	hp := tinykv.NewHTTPPool(addr)
	hp.Set(peers...)
	kv.RegisterPeers(hp)
	log.Println("tinykv is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], hp))
}

func main() {

	var configFile string
	flag.StringVar(&configFile, "c", "", "config file")
	flag.Parse()

	if err := config.LoadConfig(configFile); err != nil {
		panic(err)
	}

	group := createGroup()
	startCacheServer(config.Addr(), []string(config.Peers()), group)
}
