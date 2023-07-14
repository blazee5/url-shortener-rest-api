package models

type ShortUrl struct {
	ID  string `bson:"_id"`
	URL string `bson:"url"`
}
