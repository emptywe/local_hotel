package server

import "net/http"

// Server custom server struct
type Server struct {
	http    *http.Server
	Addr    string
	Handler http.Handler
}

// Run start listening and reving http server
func (s *Server) Run() error {
	s.http = &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}
	return s.http.ListenAndServe()
}
