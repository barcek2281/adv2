package server

import (
	"fmt"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/handler"
)

type Server struct {
	config  *config.Config
	mux     *http.ServeMux
	handlerComics *handler.HandlerComics
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config:  config,
		mux:     http.NewServeMux(),
		handlerComics: handler.NewHandler(config),
	}
}

func (s *Server) Start() error {
	s.Configure()
	return http.ListenAndServe(fmt.Sprintf(":%s", s.config.Addr), s.mux)
}

func (s *Server) Configure() {
	s.mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world"))
	})
	s.mux.Handle("POST /product/comics", s.handlerComics.Create())
	s.mux.Handle("GET /product/comics/{id}", s.handlerComics.GetById())
	s.mux.Handle("PATCH /product/comics/{id}", s.handlerComics.Update())
	s.mux.Handle("DELETE /product/comics/{id}", s.handlerComics.Delete())
	s.mux.Handle("GET /product/comics", s.handlerComics.GetAll())

}
