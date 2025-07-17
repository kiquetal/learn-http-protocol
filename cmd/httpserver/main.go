package main

import (
	"github.com/kiquetal/learn-http-protocol/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	serv, err := server.Serve(port)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	defer serv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
