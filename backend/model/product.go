package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductComics struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	Price       float64            `bson:"price" json:"price"`
	ReleaseDate string             `bson:"release_date" json:"release_date"` // or use time.Time if you prefer
	Description string             `bson:"description" json:"description"`
}
