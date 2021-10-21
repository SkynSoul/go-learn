package gohttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const ServerMetaKey = "ServerMeta"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type HttpHandler struct {
	ServerMeta map[string]string
	Handlers map[string]map[string]HandlerFunc
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawCtx := r.Context()
	valueCtx := context.WithValue(rawCtx, ServerMetaKey, h.ServerMeta)
	req := r.WithContext(valueCtx)

	method := r.Method
	urlPath := r.URL.Path
	log.Println(fmt.Sprintf("method: %s, path: %s", method, urlPath))

	if handlerMap, ok := h.Handlers[method]; ok {
		if handlerFunc, ok := handlerMap[urlPath]; ok {
			handlerFunc(w, req)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("404 page not found"))
	if err != nil {
		log.Fatal(fmt.Sprintf("method: %s, path: %s, err: %v", method, urlPath, err))
	}
}