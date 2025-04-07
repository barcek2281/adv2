package main

import (
	"log"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/server"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
}

func main() {
	c := config.NewConfig()
	s := server.NewServer(c)
	log.Printf("server start on port: %v", c.Addr)
	if err := s.Start(); err != nil {
		log.Printf("error with server: %v", err)
	}
}
