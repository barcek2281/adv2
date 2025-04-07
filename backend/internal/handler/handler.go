package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/service"
	models "github.com/barcek2281/adv2/model"
	"github.com/barcek2281/adv2/utils"
)

type HandlerComics struct {
	config      *config.Config
	serveComics *service.ComicsService
}

func NewHandler(config *config.Config) *HandlerComics {
	serveComics := service.NewComicsService(config)
	return &HandlerComics{
		config:      config,
		serveComics: serveComics,
	}
}

func (h *HandlerComics) Create() http.HandlerFunc {
	type Response struct {
		MSG string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var reqComics models.ProductComics
		if err := json.NewDecoder(r.Body).Decode(&reqComics); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			log.Printf("cannot parse, %v", err)
			return
		}
		err := h.serveComics.CreateComics(&reqComics)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			log.Printf("cannot create, %v", err)
			return
		}

		utils.Response(w, r, http.StatusCreated, Response{MSG: "comics succesfully created"})
		log.Printf("create succesfully comics")
	}
}

func (h *HandlerComics) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			log.Printf("no id")
			utils.Error(w, r, http.StatusBadRequest, nil)
			return
		}

		m, err := h.serveComics.GetById(id)
		if err != nil{
			log.Printf("repo problem: %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusAccepted, m)
	}
}