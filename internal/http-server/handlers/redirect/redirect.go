package redirect

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

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=URLGetter
type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Response{
				Status: "Error",
				Error:  "invalid request",
			})

			return
		}

		resURL, err := urlGetter.GetURL(context.Background(), alias)
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("failed to get url", "alias", alias)

			render.JSON(w, r, response.Response{
				Status: "Error",
				Error:  "not found",
			})

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Response{
				Status: "Error",
				Error:  "internal error",
			})

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
