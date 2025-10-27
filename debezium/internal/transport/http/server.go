package http

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(port int) *Server {
	srv := http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &Server{
		srv: &srv,
	}
}

func (s *Server) RegisterHandlers() error {
	http.HandleFunc("GET /api/v1/users", func(w http.ResponseWriter, r *http.Request) {})
	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
