package delete

import (
	"context"
	"errors"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/api/response"
	sl "github.com/blazee5/url-shortener-rest-api/internal/lib/logger/slog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=URLDeleter
func Delete(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Response{
				Status: "Error",
				Error:  "alias is empty",
			})

			return
		}

		err := urlDeleter.DeleteURL(context.Background(), alias)
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("failed to delete url", "alias", alias)

			render.JSON(w, r, response.Response{
				Status: "Error",
				Error:  "not found",
			})

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, response.Response{
				Status: "Error",
				Error:  "internal error",
			})

			return
		}

		log.Info("delete url", slog.String("url", alias))

		render.JSON(w, r, response.Response{
			Status: "success",
		})
	}
}
