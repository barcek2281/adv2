package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/service"
	models "github.com/barcek2281/adv2/model"
	"github.com/barcek2281/adv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if err != nil {
			log.Printf("repo problem: %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusAccepted, m)
	}
}

func (h *HandlerComics) Update() http.HandlerFunc {
	type Response struct {
		MSG string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			log.Printf("no id")
			utils.Error(w, r, http.StatusBadRequest, nil)
			return
		}

		var reqComics models.ProductComics
		if err := json.NewDecoder(r.Body).Decode(&reqComics); err != nil {
			log.Printf("not updated, %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		newId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("no id %v, %s", err, id)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		reqComics.ID = newId
		err = h.serveComics.Update(&reqComics)
		if err != nil {
			log.Printf("not updated, %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		log.Printf("updated")
		utils.Response(w, r, http.StatusResetContent, Response{MSG: "comics changed"})
	}
}

func (h *HandlerComics) Delete() http.HandlerFunc {
	type Response struct {
		MSG string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			log.Printf("no id")
			utils.Error(w, r, http.StatusBadRequest, nil)
			return
		}

		log.Printf("delete successfully")
		h.serveComics.Delete(id)
		utils.Response(w, r, http.StatusAccepted, Response{MSG: "delete successfully"})
	}
}

func (h *HandlerComics) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := h.serveComics.GetAll()
		if err != nil {
			log.Printf("get all %v", err)
			return
		}
		utils.Response(w, r, http.StatusAccepted, m)
	}
}
