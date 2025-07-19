package main

import (
	"github.com/kiquetal/learn-http-protocol/internal/request"
	"github.com/kiquetal/learn-http-protocol/internal/server"
	"github.com/kiquetal/learn-http-protocol/internal/utils"
	"io"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	// Initialize logger with INFO level
	utils.InitLogger(utils.LogLevelDebug)

	serv, err := server.Serve(port, createCustomHandler())
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

func createCustomHandler() server.Handler {
	return func(w io.Writer, rq *request.Request) *server.HandlerError {
		methodAndPath := rq.RequestLine.Method + " " + rq.RequestLine.RequestTarget
		utils.Logger.Debug("Handling request: %s", methodAndPath)
		switch methodAndPath {
		case "GET /yourproblem":
			return &server.HandlerError{
				StatusCode: 400,
				Message:    "Your problem is not my problem",
			}
		case "GET /myproblem":
			return &server.HandlerError{
				StatusCode: 500,
				Message:    "Woopsie, my bad",
			}
		case "GET /use-nvim":
			_, _ = w.Write([]byte("All good, frfr"))
			return nil

		default:
			return &server.HandlerError{
				StatusCode: 404,
				Message:    "Not Found",
			}
		}
	}
}
