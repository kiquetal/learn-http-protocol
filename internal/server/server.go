package server

import (
	"github.com/kiquetal/learn-http-protocol/internal/utils"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Port   int
	Server net.Listener
	Closed atomic.Bool // Use atomic for thread-safe boolean
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		Port:   port,
		Server: listener,
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

	buf := make([]byte, 4096) // Buffer to read data
	n, err := conn.Read(buf)
	if err != nil {
		utils.Logger.Error("Error reading from connection: %v", err)
		return
	}
	utils.Logger.Debug("Received %d bytes: %s", n, string(buf[:n]))

	// implement a response
	body := "Hello World!"

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " +
		strconv.Itoa(len(body)) + "\r\n\r\n" + body

	_, err = conn.Write([]byte(response))
	if err != nil {
		utils.Logger.Error("Error writing to connection: %v", err)
		return
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
