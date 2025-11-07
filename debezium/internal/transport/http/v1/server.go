package v1

import (
	"context"
	"debezium_server/internal/repository"
	"debezium_server/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	defaultHeaderTimeout = time.Second * 5
)

type Server struct {
	srv *http.Server
	db  *pgxpool.Pool
}

func NewServer(port int, db *pgxpool.Pool) *Server {
	srv := http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           nil,
		ReadHeaderTimeout: defaultHeaderTimeout,
	}
	return &Server{
		srv: &srv,
		db:  db,
	}
}

func (s *Server) RegisterHandlers() error {
	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo)
	handler := NewHandlerFacade(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetUsers(w, r)
	})

	s.srv.Handler = LoggingMiddleware()(mux)

	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
