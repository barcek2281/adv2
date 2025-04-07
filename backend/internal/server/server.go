package server

import (
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/handler"
)

type Server struct {
	config  *config.Config
	mux     *http.ServeMux
	handler *handler.Handler
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.Addr, s.mux)
}

func (s *Server) Configure() {
	s.mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	})
}
