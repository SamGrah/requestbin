package app

import (
	"net/http"
)

type HttpServer struct {
	httpServer *http.Server
}

func NewServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (hs *HttpServer) Start() error {
	err := hs.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
