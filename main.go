package main

import (
	"Diplom/back/db"
	"Diplom/back/server"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		slog.Info(".env файл не найден, используются дефолтные настройки")
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("db init error:%v", err)
	}
	defer db.DB.Close()
	srv := server.NewServer()

	if err := srv.Start(); err != nil {
		log.Fatalf("server  fatal error:%v", err)
	}

}
