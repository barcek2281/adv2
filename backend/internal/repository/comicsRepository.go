package repository

import "go.mongodb.org/mongo-driver/mongo"

type ComicsRepository struct {
	Collection *mongo.Collection
}
