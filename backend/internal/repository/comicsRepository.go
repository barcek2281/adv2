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

func (c *ComicsRepository) Update(comics *models.ProductComics) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"title":        comics.Title,
			"author":       comics.Author,
			"price":        comics.Price,
			"release_date": comics.ReleaseDate,
			"description":  comics.Description,
		},
	}

	result := c.Collection.FindOneAndUpdate(ctx, bson.M{"_id": comics.ID}, update)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (c *ComicsRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.Collection.DeleteOne(ctx, bson.M{"_id": newId})
	if err != nil {
		return err
	}

	return nil
}

func (c *ComicsRepository) GetAll() ([]models.ProductComics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := c.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comicsList []models.ProductComics
	for cursor.Next(ctx) {
		var comic models.ProductComics
		if err := cursor.Decode(&comic); err != nil {
			return nil, err
		}
		comicsList = append(comicsList, comic)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return comicsList, nil
}
