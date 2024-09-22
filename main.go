package main

import (
	"log"
	"net/http"
	"sync"
	"time"

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
	ss := protocol.NewBluebellServer("tcp", "0.0.0.0:9000", true)
	options := []gnet.Option{
		gnet.WithMulticore(true),               // 启用多核模式
		gnet.WithReusePort(true),               // 启用端口重用
		gnet.WithTCPKeepAlive(time.Minute * 5), // 启用 TCP keep-alive
		gnet.WithReadBufferCap(2048 * 1024),
		gnet.WithWriteBufferCap(2048 * 1024),
	}
	err := gnet.Run(ss, ss.Network+"://"+ss.Addr, options...)
	logging.Infof("server exits with error: %v", err)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go NewTCPPool(&wg)
	wg.Wait()
}
