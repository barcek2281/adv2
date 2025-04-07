package service

import (
	"log"

	"github.com/barcek2281/adv2/internal/config"
	"github.com/barcek2281/adv2/internal/store"
	models "github.com/barcek2281/adv2/model"
)

type ComicsService struct {
	config *config.Config
	store  *store.Store
}

func NewComicsService(config *config.Config) *ComicsService {
	store, err := store.NewStore(config)
	if err != nil {
		log.Fatalf("lol cannot connect to db: %v", err)
	}
	return &ComicsService{
		config: config,
		store:  store,
	}
}

func (c *ComicsService) CreateComics(comics *models.ProductComics) error {

	if err := c.store.ComicsRepo().Create(comics); err != nil {
		return err
	}

	return nil
}