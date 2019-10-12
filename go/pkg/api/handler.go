package api

import (
	"github.com/golang/protobuf/ptypes/empty"
)

// Server represents the gRPC server
type Server struct {
	Temperature float64
	Stream      TemperatureService_GetTemperatureServer
}

// SayHello generates response to a Ping request
func (s *Server) GetTemperature(e *empty.Empty, stream TemperatureService_GetTemperatureServer) error {
	s.Stream = stream

	select {}
	return nil
}
