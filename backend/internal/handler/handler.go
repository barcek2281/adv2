package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/service"
	models "github.com/barcek2281/adv2/model"
)

type HandlerComics struct {
	config  *config.Config
	serveComics *service.ComicsService
}

func NewHandler(config *config.Config) *HandlerComics {
	serveComics := service.NewComicsService(config)
	return &HandlerComics{
		config: config,
		serveComics: serveComics,
	}
}

func (h *HandlerComics) Create() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var reqComics models.ProductComics
		if err := json.NewDecoder(r.Body).Decode(&reqComics); err != nil {
			log.Printf("cannot parse, %v", err)
			return
		}
		err := h.serveComics.CreateComics(&reqComics)
		if err != nil {
			log.Printf("cannot create, %v", err)
			return
		}


		log.Printf("create succesfully comics")
	}
}
