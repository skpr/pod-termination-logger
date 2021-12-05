package server

import "io"

// Server for handling Pod events.
type Server struct {
	Writer io.Writer
}

// New server for handling Pod events.
func New(w io.Writer) *Server {
	return &Server{Writer: w}
}
