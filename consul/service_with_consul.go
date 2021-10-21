package consul

import (
	"fmt"
	gohttp "github.com/SkynSoul/go-learn/go-http"
	consulApi "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ServerWithConsul struct {
	wg   	   	*sync.WaitGroup
	Host   	   	string
	Port       	int
	Server     	*gohttp.SimpleHttpServer
	AgentProxy 	*AgentProxy
	Name       	string
	ConsulAddr 	string
	ConsulDC   	string
}

func NewService(wg *sync.WaitGroup, name string, consulAddr string, consulDC string, host string, port int) (*ServerWithConsul, error) {
	s := &ServerWithConsul{}
	s.wg = wg
	s.Host = host
	s.Port = port
	s.Name = name
	s.ConsulAddr = consulAddr
	s.ConsulDC = consulDC

	var err error

	s.Server, err = gohttp.NewHttpServer(wg, host, port)
	if err != nil {
		return nil, err
	}
	s.registerHealthHandler()

	s.AgentProxy, err = NewAgentProxy(s.ConsulAddr, s.ConsulDC)
	if err != nil {
		return nil, err
	}

	return s, err
}

func (s *ServerWithConsul) Start(isRegister bool) {
	if s.Server == nil {
		return
	}
	s.wg.Add(1)

	go s.Server.Start()

	go s.catchSignal()

	if isRegister {
		s.registerService()
	}
}

func (s *ServerWithConsul) registerHealthHandler() {
	s.Server.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			log.Println("health check write filed, err: ", err)
		}
	})
}

func (s *ServerWithConsul) catchSignal() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	s.Stop()
}

func (s *ServerWithConsul) registerService() {
	tags := []string{
		"consul.native=true",
	}
	check := &consulApi.AgentServiceCheck{
		HTTP: fmt.Sprintf("http://%s:%d/health", s.Host, s.Port),
		Interval: "10s",
		Timeout: "2s",
	}
	err := s.AgentProxy.RegisterService(s.Name, s.Host, s.Port, tags, check)
	if err != nil {
		log.Println("Server register failed: ", err)
		s.Stop()
	}
	log.Printf("Server register success")
}

func (s *ServerWithConsul) Stop() {
	defer s.wg.Done()

	if s.Server == nil {
		return
	}

	err := s.AgentProxy.DeregisterService(s.Name, s.Host, s.Port)
	if err != nil {
		log.Println("Server deregister failed: ", err)
	}
	log.Printf("Server deregister success")

	s.Server.Stop()
}