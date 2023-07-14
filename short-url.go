package url_shortener_rest

type ShortUrl struct {
	URL   string `bson:"url"`
	Alias string `bson:"alias"`
}
