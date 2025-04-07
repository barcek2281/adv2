package repository

import (
	"context"
	"fmt"
	"time"

	models "github.com/barcek2281/adv2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComicsRepository struct {
	Collection *mongo.Collection
}

func (c *ComicsRepository) Create(comics *models.ProductComics) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := c.Collection.FindOne(ctx, bson.M{"_id": comics.ID}).Decode(comics)

	if err == nil {
		return fmt.Errorf("comis with id exist already")
	}

	_, err = c.Collection.InsertOne(ctx, comics)
	if err != nil {
		return err
	}

	return nil
}

func (c *ComicsRepository) GetById(id string) (*models.ProductComics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	var comics models.ProductComics
	err = c.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&comics)
	if err != nil {
		return nil, err
	}

	return &comics, nil
}
