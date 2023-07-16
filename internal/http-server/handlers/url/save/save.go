package save

import (
	"context"
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

type URLSaver interface {
	SaveUrl(ctx context.Context, shortUrl *models.ShortUrl) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		err = urlSaver.SaveUrl(context.Background(), &models.ShortUrl{
			ID:  alias,
			URL: req.URL,
		})
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
