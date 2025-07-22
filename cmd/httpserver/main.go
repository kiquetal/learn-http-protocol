package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/kiquetal/learn-http-protocol/internal/headers"
	"github.com/kiquetal/learn-http-protocol/internal/request"
	"github.com/kiquetal/learn-http-protocol/internal/response"
	"github.com/kiquetal/learn-http-protocol/internal/server"
	"github.com/kiquetal/learn-http-protocol/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

		default:
			//check path begins with /httpbin
			if strings.HasPrefix(rq.RequestLine.RequestTarget, "/httpbin") {
				utils.Logger.Info("Handling /httpbin request")
				subPath := strings.TrimPrefix(rq.RequestLine.RequestTarget, "/httpbin")
				utils.Logger.Debug("Handling /httpbin request: %s", subPath)
				res, err := http.Get("https://httpbin.org" + subPath) // Simulate a request to httpbin.org
				if err != nil {
					utils.Logger.Error("Error fetching from httpbin.org: %v", err)
					_ = w.WriteStatusLine(500) // HTTP 500 Internal Server Error
					_ = w.WriteHeaders(response.GetDefaultHeaders(len(getInternalServerErrorHtml())))
					_, _ = w.Write([]byte(getInternalServerErrorHtml()))
					return
				}
				defer res.Body.Close()
				n := make([]byte, 512) // Initialize a byte slice with capacity
				// read all the body

				_ = w.WriteStatusLine(200) // HTTP 200 OK
				//create headers for a response

				headersForResponse := headers.NewHeaders()
				headersForResponse["Content-Type"] = res.Header.Get("Content-Type")
				headersForResponse["Transfer-Encoding"] = "chunked"
				headersForResponse["Trailer"] = "X-Content-SHA256,X-Content-Length"
				_ = w.WriteHeaders(headersForResponse)
				hash := sha256.New()
				totalBytes := 0
				for {
					readbytes, err := res.Body.Read(n)
					if err != nil {
						if err.Error() == "EOF" {
							//	intbytes, err := w.WriteChunkedBodyDone()
							hashFinal := hash.Sum(nil)
							headerForHash := headers.Header{
								"X-Content-SHA256": fmt.Sprintf("%x", hashFinal),
								"X-Content-Length": strconv.Itoa(totalBytes),
							}
							w.WriteTrailers(headerForHash)

							if err != nil {
								utils.Logger.Error("Error writing chunked body done: %v", err)
								return
							}
							break // End of file, stop reading
						}
						utils.Logger.Error("Error reading response body: %v", err)
						_ = w.WriteStatusLine(500) // HTTP 500 Internal Server Error
						_ = w.WriteHeaders(response.GetDefaultHeaders(len(getInternalServerErrorHtml())))
						_, _ = w.Write([]byte(getInternalServerErrorHtml()))
						return

					}
					utils.Logger.Info("Response body: %s", string(n[:readbytes]))
					//return chunked data using the function WriteChunked
					if readbytes > 0 {
						hash.Write(n[:readbytes])
						totalBytes += readbytes
						writeChunkedBody, err := w.WriteChunkedBody(n[:readbytes])
						if err != nil {
							return

						}
						if writeChunkedBody > 0 {
							utils.Logger.Debug("Wrote chunked data: %s", string(n[:readbytes]))
						}
					}
				}

			} else {

				_ = w.WriteStatusLine(404) // HTTP 404 Not Found
				_ = w.WriteHeaders(response.GetDefaultHeaders(len(getNotFoundHtml())))
				_, _ = w.Write([]byte(getNotFoundHtml()))
			}
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
