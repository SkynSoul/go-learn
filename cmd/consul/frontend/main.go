package main

import (
	"flag"
	"github.com/SkynSoul/go-learn/consul"
	"github.com/SkynSoul/go-learn/utils"
	"log"
	"sync"
)

var host = flag.String("host", "", "Service Host")
var port = flag.Int("port", 80, "Service Port")
var name = flag.String("name", "go-server-frontend", "Service Name")
var consulAddr = flag.String("consuladdr", "127.0.0.1:8500", "Consul Address")
var consulDC = flag.String("dc", "local-vm", "Consul Datacenter")

func main() {
	flag.Parse()
	if *host == "" {
		*host = utils.GetLocalIP()[0]
	}

	wg := &sync.WaitGroup{}
	s, err := consul.NewService(wg, *name, *consulAddr, *consulDC, *host, *port)
	if err != nil {
		log.Fatalln("create http server failed: ", err)
	}
	consul.RegisterFrontHandler(s, "go-server-backend")
	s.Start(true)
	wg.Wait()
}