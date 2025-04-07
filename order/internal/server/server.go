package server

import (
	"fmt"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/handler"
)

type Server struct {
	config        *config.Config
	mux           *http.ServeMux
	handlerOrders *handler.HandlerOrders
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config:        config,
		mux:           http.NewServeMux(),
		handlerOrders: handler.NewHandlerOrders(config),
	}
}

func (s *Server) Start() error {
	s.Configure()
	return http.ListenAndServe(fmt.Sprintf(":%s", s.config.Addr), s.mux)
}

func (s *Server) Configure() {
	// Пинг
	s.mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	})

	// Orders endpoints
	s.mux.Handle("POST /orders", s.handlerOrders.Create())
	s.mux.Handle("GET /orders/{id}", s.handlerOrders.GetByID())
	s.mux.Handle("PATCH /orders/{id}", s.handlerOrders.Update())
	s.mux.Handle("DELETE /orders/{id}", s.handlerOrders.Delete())
	s.mux.Handle("GET /orders", s.handlerOrders.GetAll())
}
