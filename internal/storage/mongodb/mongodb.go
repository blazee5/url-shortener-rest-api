package mongodb

import (
	"context"
	"fmt"
	models "github.com/blazee5/url-shortener-rest-api"
	"github.com/blazee5/url-shortener-rest-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
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

	return &Storage{DB: client}, nil
}

type UrlDAO struct {
	c *mongo.Collection
}

func NewDAO(client *mongo.Client) (*UrlDAO, error) {
	return &UrlDAO{
		c: client.Database("url-shortener").Collection("shortUrls"),
	}, nil
}

func (dao *UrlDAO) SaveURL(ctx context.Context, urlToSave string, alias string) error {
	_, err := dao.c.InsertOne(ctx, models.ShortUrl{
		ID:  alias,
		URL: urlToSave,
	})

	if err != nil {
		return err
	}

	return nil
}

func (dao *UrlDAO) GetURL(ctx context.Context, alias string) (string, error) {
	filter := bson.D{{"_id", alias}}
	var URL models.ShortUrl
	err := dao.c.FindOne(ctx, filter).Decode(&URL)
	if err != nil {
		return URL.URL, err
	}

	return URL.URL, nil
}

func (dao *UrlDAO) DeleteURL(ctx context.Context, alias string) error {
	filter := bson.D{{"_id", alias}}
	_, err := dao.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
