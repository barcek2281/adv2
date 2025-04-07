package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/service"
	model "github.com/barcek2281/adv2/models"
	"github.com/barcek2281/adv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HandlerOrders struct {
	config        *config.Config
	ordersService *service.OrdersService
}

func NewHandlerOrders(config *config.Config) *HandlerOrders {
	ordersService := service.NewOrdersService(config)
	return &HandlerOrders{
		config:        config,
		ordersService: ordersService,
	}
}

func (h *HandlerOrders) Create() http.HandlerFunc {
	type Response struct {
		MSG string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var reqOrder model.Order
		if err := json.NewDecoder(r.Body).Decode(&reqOrder); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			log.Printf("cannot parse, %v", err)
			return
		}
		err := h.ordersService.CreateOrder(&reqOrder)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			log.Printf("cannot create, %v", err)
			return
		}

		utils.Response(w, r, http.StatusCreated, Response{MSG: "order successfully created"})
		log.Printf("create successfully order")
	}
}

func (h *HandlerOrders) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			log.Printf("no id")
			utils.Error(w, r, http.StatusBadRequest, nil)
			return
		}

		order, err := h.ordersService.GetOrderByID(id)
		if err != nil {
			log.Printf("repo problem: %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusAccepted, order)
	}
}

func (h *HandlerOrders) Update() http.HandlerFunc {
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

		var reqOrder model.Order
		if err := json.NewDecoder(r.Body).Decode(&reqOrder); err != nil {
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

		reqOrder.ID = newId
		err = h.ordersService.UpdateOrder(&reqOrder)
		if err != nil {
			log.Printf("not updated, %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		log.Printf("updated")
		utils.Response(w, r, http.StatusResetContent, Response{MSG: "order changed"})
	}
}

func (h *HandlerOrders) Delete() http.HandlerFunc {
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
		err := h.ordersService.DeleteOrder(id)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusAccepted, Response{MSG: "delete successfully"})
	}
}

func (h *HandlerOrders) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := h.ordersService.GetAllOrders()
		if err != nil {
			log.Printf("get all %v", err)
			return
		}
		utils.Response(w, r, http.StatusAccepted, orders)
	}
}
