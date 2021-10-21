package consul

import (
	"fmt"
	gohttp "github.com/SkynSoul/go-learn/go-http"
	"log"
	"net/http"
	"path"
	"text/template"
)

func WrapHandlerFrontendRoot(s *ServerWithConsul, serviceName string) gohttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceMap := s.AgentProxy.GetService(serviceName)
		addrList := make([]string, len(serviceMap))
		idx := 0
		for _, info := range serviceMap {
			addrList[idx] = fmt.Sprintf("http://%s:%d", info.Host, info.Port)
			idx++
		}
		log.Println(addrList)
		tmpl, err := template.ParseFiles(path.Join(WorkingPath, "./static/html/lb.html"))
		if err != nil {
			log.Println(fmt.Sprintf("%s process error, err: %v", r.URL.Path, err))
			HandlerError(w, r)
			return
		}
		tmpl.Execute(w, addrList)
	}
}

func RegisterFrontHandler(s *ServerWithConsul, discoverName string) {
	s.Server.Get("/", WrapHandlerFrontendRoot(s, discoverName))
}
