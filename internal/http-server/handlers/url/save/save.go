package save

import (
	"context"
	"errors"
	models "github.com/blazee5/url-shortener-rest-api"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/api/response"
	sl "github.com/blazee5/url-shortener-rest-api/internal/lib/logger/slog"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/random"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
)

// TODO: move to config
const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=URLSaver
type URLSaver interface {
	SaveURL(ctx context.Context, shortUrl *models.ShortUrl) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		for i := 1; i <= 10; i++ {
			err = urlSaver.SaveURL(context.Background(), &models.ShortUrl{
				ID:  alias,
				URL: req.URL,
			})

			if err == nil {
				break
			}

			alias = random.NewRandomString(aliasLength)
		}

		if errors.Is(err, errors.New("url exists")) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "url already exists",
			})

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to add url",
			})

			return
		}

		render.JSON(w, r, Response{
			Status: "success",
			Alias:  alias,
		})
	}
}
