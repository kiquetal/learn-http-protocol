package main

import (
	"fmt"
	"github.com/kiquetal/learn-http-protocol/internal/server"
	"github.com/kiquetal/learn-http-protocol/internal/utils"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	// Initialize logger with INFO level
	utils.InitLogger(utils.LogLevelInfo)

	serv, err := server.Serve(port)
	if err != nil {
		utils.Logger.Debug("Failed to start server: %v", err)
		utils.Logger.Error("Failed to serve: %v", err)
		os.Exit(1)
	}
	defer serv.Close()

	utils.Logger.Info("Server started on port %d", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	utils.Logger.Info("Server gracefully stopped")
}

type Handler func(conn net.Conn)
type HandlerError struct {
	StatusCode int
	Message    string
}

func WriteErrorResponse(w io.Writer, handleErr HandlerError) {

	_, err := fmt.Fprint(w, "HTTP/1.1 ", handleErr.StatusCode, " ", handleErr.Message, "\r\n")
	if err != nil {
		utils.Logger.Error("Error writing response: %v", err)
	}
}
