package store

import (
	"context"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	uri        string
	config     *config.Config
	db         *mongo.Database
	client     *mongo.Client
	comicsRepo *repository.ComicsRepository
}

func NewStore(config *config.Config) (*Store, error) {
	ClientOptions := options.Client().ApplyURI(config.Uri)
	client, err := mongo.Connect(context.TODO(), ClientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(config.DBname)

	return &Store{
		uri:    config.Uri,
		db:     db,
		client: client,
		config: config,
	}, nil
}

func (s *Store) ComicsRepo() *repository.ComicsRepository {
	if s.comicsRepo == nil {
		s.comicsRepo = &repository.ComicsRepository{Collection: s.db.Collection(s.config.ComicsCollection)}
	}
	return s.comicsRepo
}
