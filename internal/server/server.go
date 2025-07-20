package server

import (
	"bytes"
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
			s.handle(c)
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

	utils.Logger.Debug("Starting to handle connections")
	req, e := request.RequestFromReader(conn)
	if e != nil {
		utils.Logger.Error("Error parsing request: %v", e)
		WriteErrorResponse(conn, HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Bad Request",
		})
		return
	}
	utils.Logger.Debug("Parsed request: %v", req)
	writer := response.Writer{
		Writer:      conn,
		WriteStatus: response.WriterStatusInitialized,
	}
	bufferForHandler := new(bytes.Buffer)
	s.handler(writer, req)
	if handleErr != nil {
		utils.Logger.Error("Handler error: %v", handleErr)
		WriteErrorResponse(conn, *handleErr)
		return
	}
	responseBuffer := new(bytes.Buffer)
	resp := response.Response{StatusCode: response.StatusOK}
	headers := response.GetDefaultHeaders(len(bufferForHandler.Bytes()))
	response.WriteStatusLine(responseBuffer, resp.StatusCode)
	response.WriteHeaders(responseBuffer, headers)
	responseBuffer.Write(bufferForHandler.Bytes())
	responseBuffer.WriteTo(conn)
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

type Handler func(r *response.Writer, req *request.Request)
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func WriteErrorResponse(w io.Writer, handleErr HandlerError) {

	err := response.WriteStatusLine(w, handleErr.StatusCode)
	if err != nil {
		return
	}
	headers := response.GetDefaultHeaders(len(handleErr.Message))

	err = response.WriteHeaders(w, headers)
	if err != nil {
		utils.Logger.Error("Error writing headers: %v", err)
		return
	}
	_, err = w.Write([]byte(handleErr.Message))

}
