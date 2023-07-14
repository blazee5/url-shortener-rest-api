package mongodb

import (
	"context"
	"fmt"
	"github.com/blazee5/url-shortener-rest-api/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	DB *mongo.Client
}

func Run(cfg *config.Config) (*Storage, error) {
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort))
	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	return &Storage{DB: client}, nil
}

type UrlDAO struct {
	c *mongo.Collection
}

func NewDAO(ctx context.Context, client *mongo.Client) (*UrlDAO, error) {
	return &UrlDAO{
		c: client.Database("url-shortener").Collection("shortUrls"),
	}, nil
}
