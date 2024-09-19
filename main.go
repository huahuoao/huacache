package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	huacache "github.com/huahuoao/huacache/core"
	"github.com/huahuoao/huacache/core/protocol"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
)

func NewHTTPPool(wg *sync.WaitGroup) {
	defer wg.Done()
	addr := "0.0.0.0:4160"
	peers := huacache.NewHTTPPool(addr)
	log.Println("gcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

func NewTCPPool(wg *sync.WaitGroup) {
	defer wg.Done()
	var port int
	var multicore bool
	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore=true")
	flag.Parse()
	ss := protocol.NewBluebellServer("tcp", "localhost:9000", true)
	err := gnet.Run(ss, ss.Network+"://"+ss.Addr, gnet.WithMulticore(multicore))
	logging.Infof("server exits with error: %v", err)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // 等待两个 goroutine
	//	go NewHTTPPool(&wg)
	go NewTCPPool(&wg)
	wg.Wait() // 等待所有 goroutine 完成
}
