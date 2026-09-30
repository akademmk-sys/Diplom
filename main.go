package main

import (
	"Diplom/back/db"
	"Diplom/back/server"
	"log"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("db init error:%v", err)
	}
	srv := server.NewServer()

	if err := srv.Start(); err != nil {
		log.Fatalf("server  fatal error:%v", err)
	}

}
