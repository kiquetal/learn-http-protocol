package server

import (
	"fmt"
	"github.com/kiquetal/learn-http-protocol/internal/request"
	"github.com/kiquetal/learn-http-protocol/internal/response"
	"github.com/kiquetal/learn-http-protocol/internal/utils"
	"io"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Port    int
	Server  net.Listener
	handler Handler     // Handler function to process connections
	Closed  atomic.Bool // Use atomic for thread-safe boolean
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		Port:    port,
		Server:  listener,
		handler: handler,
	}
	server.Closed.Store(false)
	go server.listen()
	return server, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Server.Accept()
		if s.Closed.Load() {
			utils.Logger.Info("Server is closed, stopping accept loop")
			return // Exit the loop if the server is closed
		}
		if err != nil {
			utils.Logger.Warn("Error accepting connection: %v", err)
			continue // Handle error appropriately in production code
		}
		go func(c net.Conn) {
			// call the handle function to process the connection
			s.handler(c)

			defer c.Close()
			// Handle the connection (e.g., read request, send response)
			// This is where you would implement your request handling logic
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	// This function should handle the connection, read the request,
	// parse it, and send a response.
	// For now, we just close the connection.
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			utils.Logger.Error("Error closing connection: %v", err)
		}
	}(conn)

	buf := make([]byte, 4096) // Buffer to read data
	n, err := conn.Read(buf)
	if err != nil {
		utils.Logger.Error("Error reading from connection: %v", err)
		return
	}
	utils.Logger.Debug("Received %d bytes: %s", n, string(buf[:n]))

	headers := response.GetDefaultHeaders(0)
	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		utils.Logger.Error("Error writing status line: %v", err)
	}
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		utils.Logger.Error("Error writing headers to response: %v", err)
	}

	utils.Logger.Info("Response sent successfully")
}

func (s *Server) Close() error {
	if s.Closed.Load() {
		utils.Logger.Info("Server is already closed")
		return nil // Server is already closed, nothing to do
	}
	s.Closed.Store(true)

	utils.Logger.Info("Server closed successfully")
	return s.Server.Close() // Close the listener
}

type Handler func(conn net.Conn) *HandlerError
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func WriteErrorResponse(w io.Writer, handleErr HandlerError) {

	_, err := fmt.Fprint(w, "HTTP/1.1 ", handleErr.StatusCode, " ", handleErr.Message, "\r\n")
	if err != nil {
		utils.Logger.Error("Error writing response: %v", err)
	}
}

func handlerFunction(conn net.Conn) *HandlerError {

	//parse the request from the connection

	r, err := request.RequestFromReader(conn)
	if err != nil {
		utils.Logger.Error("Error reading request: %v", err)
		return &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Bad Request",
		}
	}

	return nil
}
