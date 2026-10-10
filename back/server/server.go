package server

import (
	"Diplom/back/handlers"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	port   string
	router *chi.Mux
}

func NewServer() *Server {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.Init(r)
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	return &Server{
		port:   port,
		router: r,
	}

}

func (s *Server) Start() error {
	addr := ":" + s.port
	slog.Info("Сервер запущен", slog.String("addr", "http://localhost"+addr))
	return http.ListenAndServe(addr, s.router)
}
