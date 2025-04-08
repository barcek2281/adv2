package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/barcek2281/adv2/internal/config"
	model "github.com/barcek2281/adv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrdersService struct {
	config *config.Config
	db     *mongo.Database
}

func NewOrdersService(config *config.Config) *OrdersService {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	db := client.Database(config.DBname)
	return &OrdersService{config: config, db: db}
}

func (s *OrdersService) CreateOrder(order *model.Order) error {

	req, err := http.Get(fmt.Sprintf("localhost:8080/product/comics/%s", order.ProductID))
	if err != nil || req.StatusCode == http.StatusBadRequest {
		return errors.New("not found product")
	}
	collection := s.db.Collection(s.config.OrdersCollection)
	_, err = collection.InsertOne(context.Background(), order)
	return err
}

func (s *OrdersService) GetOrderByID(id string) (*model.Order, error) {
	collection := s.db.Collection(s.config.OrdersCollection)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order model.Order
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrdersService) UpdateOrder(order *model.Order) error {
	collection := s.db.Collection(s.config.OrdersCollection)
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": order.ID}, bson.M{"$set": order})
	return err
}

func (s *OrdersService) DeleteOrder(id string) error {
	collection := s.db.Collection(s.config.OrdersCollection)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (s *OrdersService) GetAllOrders() ([]model.Order, error) {

	collection := s.db.Collection(s.config.OrdersCollection)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var orders []model.Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
