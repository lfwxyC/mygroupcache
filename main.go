package main

import (
	"GeeCache/geecache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"cyw": "661",
	"czy": "661",
	"ccc": "690",
}

func createGroup(cacheType geecache.CacheType) *geecache.Group {
	return geecache.NewGroup("scores", geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10, cacheType)
}

func startCacheServer(self string, addrs []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(self)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at ", self)
	log.Fatal(http.ListenAndServe(self[7:], peers))
}

func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("frontend server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	var lfu bool // 淘汰策略是否使用lfu，默认为lru
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.BoolVar(&lfu, "lfu", false, "Use lfu?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	cacheType := geecache.LRU
	if lfu {
		cacheType = geecache.LFU
	}
	gee := createGroup(geecache.CacheType(cacheType))
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
