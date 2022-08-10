package server

import (
	"log"
	"net"
	"net/http"
)

type Server struct {
	server   *http.Server
	handler  Handler
	httpPort string
	httpHost string
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func NewServer(h Handler, host, port string) *Server {
	return &Server{handler: h, httpHost: host, httpPort: port}
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.httpHost, s.httpPort),
		Handler: s.handler,
	}

	log.Println("Starting listen")
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return s.server.Close()
}
