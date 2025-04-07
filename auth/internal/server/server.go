package server

import (
	"log"
	"net/http"

	"github.com/barcek2281/adv2/auth/internal/config"
	"github.com/barcek2281/adv2/auth/internal/handler"
	"github.com/barcek2281/adv2/auth/internal/store"
	"github.com/gorilla/mux"
)

type Server struct {
	handler *handler.Handler
	config  *config.Config
	Router  *mux.Router
}

func NewServer(config *config.Config) *Server {
	router := mux.NewRouter()
	storage, err := store.NewMongoDB(config)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	handler := handler.NewHandler(storage, config)
	s := Server{Router: router, config: config, handler: handler}

	return &s
}

func (s *Server) Run() error {
	s.Configure()

	return http.ListenAndServe(s.config.Addr, s.Router)
}

func (s *Server) Configure() {
	s.Router.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`PONG`))
	})).Methods("GET")
	s.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	s.Router.HandleFunc("/user/register", s.handler.UserHandler().RegisterUser()).Methods("POST")
	s.Router.HandleFunc("/user/login", s.handler.UserHandler().Login()).Methods("POST")
	s.Router.HandleFunc("/admin/register", s.handler.AdminHandler().Register()).Methods("POST")
	s.Router.HandleFunc("/admin/login", s.handler.AdminHandler().Login()).Methods("POST")
}
