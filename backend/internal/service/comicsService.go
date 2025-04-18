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

func (c *ComicsService) GetById(id string) (*models.ProductComics, error) {
	return c.store.ComicsRepo().GetById(id)
}

func (c *ComicsService) Update(comics *models.ProductComics) error {
	return c.store.ComicsRepo().Update(comics)
}

func (c *ComicsService) Delete(id string) error {
	return c.store.ComicsRepo().Delete(id)
}

func (c *ComicsService) GetAll() ([]models.ProductComics, error) {
	return c.store.ComicsRepo().GetAll()
}
