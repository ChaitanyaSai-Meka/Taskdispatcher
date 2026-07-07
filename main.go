package main

import (
	"log"

	"github.com/ChaitanyaSai-Meka/Taskdispatcher/internal/dispatcher"
	"github.com/ChaitanyaSai-Meka/Taskdispatcher/server"
)

const (
	maxWorkers = 5
	tcpAddr    = ":9000"
	httpAddr   = ":8080"
)

func main() {
	d := dispatcher.New(maxWorkers)
	go d.Run()

	if err := server.StartTCP(tcpAddr, d); err != nil {
		log.Fatal("tcp server failed to start:", err)
	}

	log.Println("HTTP server listening on", httpAddr)
	if err := server.StartHTTP(httpAddr, d); err != nil {
		log.Fatal("http server failed:", err)
	}
}
