package server

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func NewServer() {
	if err := godotenv.Load(); err != nil {
		slog.Info("using default enviroment")
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	addr := ":" + port
	slog.Info("Сервер запущен", slog.String("addr", "http://localhost"+addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
