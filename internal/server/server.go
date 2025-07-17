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

func (s *Server) Serve(port int) (*Server, error) {

	s.Port = port
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	s.Server = listener
	return s, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Server.Accept()
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

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, World!"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing to connection: %v\n", err)
		return
	}
	fmt.Println("Response sent successfully")

}
