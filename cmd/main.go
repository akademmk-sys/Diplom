package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

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

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))

	r.Handle("/*", fileServer)

	addr := ":" + port
	slog.Info("Сервер запущен", slog.String("addr", "http://localhost"+addr))

	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("Ошибка запуска сервера", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
