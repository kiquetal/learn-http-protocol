package main

import (
	"github.com/kiquetal/learn-http-protocol/internal/request"
	"github.com/kiquetal/learn-http-protocol/internal/response"
	"github.com/kiquetal/learn-http-protocol/internal/server"
	"github.com/kiquetal/learn-http-protocol/internal/utils"
	"net/http"
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
	return func(w *response.Writer, rq *request.Request) {
		methodAndPath := rq.RequestLine.Method + " " + rq.RequestLine.RequestTarget
		utils.Logger.Debug("Handling request: %s", methodAndPath)
		switch methodAndPath {
		case "GET /yourproblem":
			_ = w.WriteStatusLine(400) // HTTP 200 OK
			_ = w.WriteHeaders(response.GetDefaultHeaders(len(getBadRequestHtml())))
			_, _ = w.Write([]byte(getBadRequestHtml()))

		case "GET /myproblem":

			_ = w.WriteStatusLine(500) // HTTP 400 Bad Request
			_ = w.WriteHeaders(response.GetDefaultHeaders(len(getInternalServerErrorHtml())))
			_, _ = w.Write([]byte(getInternalServerErrorHtml()))
		case "GET /":
			_ = w.WriteStatusLine(200) // HTTP 200 OK
			_ = w.WriteHeaders(response.GetDefaultHeaders(len(getOkHtml())))
			_, _ = w.Write([]byte(getOkHtml()))

		case "GET /httpbin/stream/100":
			get, err := http.Get("https://httpbin.org/stream/100")
			if err != nil {
				return
			}
			defer get.Body.Close()
			n := []byte
			resp := get.Body.Read(n)

		default:
			_ = w.WriteStatusLine(404) // HTTP 400 Bad Request
			_ = w.WriteHeaders(response.GetDefaultHeaders(len(getNotFoundHtml())))
			_, _ = w.Write([]byte(getNotFoundHtml()))
		}
	}
}

func getOkHtml() string {
	return `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
}
func getBadRequestHtml() string {
	return `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

}
func getInternalServerErrorHtml() string {
	return `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
}

func getNotFoundHtml() string {
	return `<html>
  <head>
	<title>404 Not Found</title>
  </head>
  <body>
	<h1>Not Found</h1>
	<p>Sorry, I couldn't find what you were looking for.</p>
  </body>
}	`
}
