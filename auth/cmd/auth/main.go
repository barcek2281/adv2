package main

import (
	"log"
	"os"

	"github.com/barcek2281/adv2/auth/internal/config"
	"github.com/barcek2281/adv2/auth/internal/server"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Printf("error with config: %v", err)
		os.Exit(1)
	}

	s := server.NewServer(c)
	log.Printf("server start on port: %s", c.Addr)

	if err := s.Run(); err != nil {
		log.Printf("server err: %v", err)
		os.Exit(1)
	}
}
