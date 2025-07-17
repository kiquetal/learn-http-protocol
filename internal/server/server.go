package server

import (
	"net"
	"strconv"
)

type Server struct {
	Port   int
	Server net.Listener
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

func (s *Server) Close() error {
	if s.Server != nil {
		return s.Server.Close()
	}
	return nil
}
