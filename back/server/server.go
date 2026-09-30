package server

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {

}

func NewServer(logger *log.Logger) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	err := http.ListenAndServe("")
}
