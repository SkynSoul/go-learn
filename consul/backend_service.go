package consul

import (
	"fmt"
	gohttp "github.com/SkynSoul/go-learn/go-http"
	"log"
	"net/http"
	"path"
	"text/template"
)

func HandlerBackendRoot(w http.ResponseWriter, r *http.Request) {
	serverMeta := r.Context().Value(gohttp.ServerMetaKey)
	tmpl, err := template.ParseFiles(path.Join(WorkingPath, "./static/html/welcome.html"))
	if err != nil {
		log.Println(fmt.Sprintf("%s process error, err: %v", r.URL.Path, err))
		HandlerError(w, r)
		return
	}
	tmpl.Execute(w, serverMeta)
}

func RegisterBackendHandler(s *ServerWithConsul) {
	s.Server.Get("/", HandlerBackendRoot)
}
