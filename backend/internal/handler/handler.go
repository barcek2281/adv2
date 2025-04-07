package handler

import (
	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/service"
)

type Handler struct {
	config  *config.Config
	service *service.ComicsService
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{
		config: config,
	}
}
