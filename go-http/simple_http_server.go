package gohttp

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	ServerStatusInit	= 0
	ServerStatusRunning	= 1
	ServerStatusLeft	= 2
)

type SimpleHttpServer struct {
	host string
	port int
	handler *HttpHandler
	server *http.Server
	wg *sync.WaitGroup
	status int8
}

func NewHttpServer(wg *sync.WaitGroup, host string, port int) (*SimpleHttpServer, error) {
	s := &SimpleHttpServer{}
	err := s.init(wg, host, port)
	return s, err
}

func (s *SimpleHttpServer) init(wg *sync.WaitGroup, host string, port int) error {
	rand.Seed(time.Now().UnixNano())
	s.host = host
	s.port = port
	s.wg = wg
	s.handler = &HttpHandler{
		Handlers: make(map[string]map[string]HandlerFunc),
		ServerMeta: map[string]string{
			"Host": s.host,
			"Port": strconv.Itoa(s.port),
		},
	}
	s.server = &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", s.port),
		Handler: s.handler,
	}
	s.status = ServerStatusInit
	return nil
}

func (s *SimpleHttpServer) Start() {
	if s.server == nil || s.status != ServerStatusInit {
		return
	}
	s.status = ServerStatusRunning
	s.wg.Add(1)

	go s.startHttpServer()
}

func (s *SimpleHttpServer) startHttpServer() {
	log.Printf("server start at %s:%d\n", s.host, s.port)
	err := s.server.ListenAndServe()
	if err != nil && err != io.EOF {
		log.Println("server start failed: ", err)
		s.Stop()
	}
}

func (s *SimpleHttpServer) Stop() {
	defer s.wg.Done()

	if s.status == ServerStatusLeft {
		return
	}
	s.status = ServerStatusLeft

	if s.server == nil {
		return
	}

	err := s.server.Close()
	if err != nil {
		log.Println("server close failed: ", err)
	} else {
		log.Println("server close success")
	}
}

func (s *SimpleHttpServer) Get(path string, f HandlerFunc) {
	if f == nil || path == "" {
		return
	}
	getMap, ok := s.handler.Handlers["GET"]
	if !ok {
		getMap = make(map[string]HandlerFunc)
		s.handler.Handlers["GET"] = getMap
	}
	getMap[path] = f
}

func (s *SimpleHttpServer) Post(path string, f HandlerFunc) {
	if f == nil || path == "" {
		return
	}
	getMap, ok := s.handler.Handlers["POST"]
	if !ok {
		getMap = make(map[string]HandlerFunc)
		s.handler.Handlers["POST"] = getMap
	}
	getMap[path] = f
}
