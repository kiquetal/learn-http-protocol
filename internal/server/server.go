package server

import (
	"fmt"
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
			fmt.Println("Server is closed, stopping accept loop")
			return // Exit the loop if the server is closed
		}
		if err != nil {
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

		}
	}(conn)

	buf := make([]byte, 4096) // Buffer to read data
	n, err := conn.Read(buf)
	if err != nil {
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
		}
	}
	fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
	// implement a response

	body := "Hello World!"

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing to connection: %v\n", err)
		return
	}
	fmt.Println("Response sent successfully")

}

func (s *Server) Close() error {
	if s.Closed.Load() {
		fmt.Println("Server is already closed")
		return nil // Server is already closed, nothing to do
	}
	s.Closed.Store(true)

	fmt.Println("Server closed successfully")
	return s.Server.Close() // Close the listener
}
