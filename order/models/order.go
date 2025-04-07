package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	ProductID string             `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Total     float64            `json:"total" bson:"total"`
	Status    string             `json:"status" bson:"status"`
}
