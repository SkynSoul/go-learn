package consul

import (
	"sync"
	"testing"
)

func TestSimpleHttpServer(t *testing.T) {
	wg := &sync.WaitGroup{}
	s, _ := NewService(wg, "go-Server-backend", "127.0.0.1:8500", "local-vm", "127.0.0.1", 1234)
	s.Start()
	wg.Wait()
}
