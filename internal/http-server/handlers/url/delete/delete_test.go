package delete

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/url/delete/mocks"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/api/response"
	sl "github.com/blazee5/url-shortener-rest-api/internal/lib/logger/slog"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "With alias",
			alias: "Dj4fKS2",
		},
		{
			name:      "Without alias",
			alias:     "",
			respError: "alias is empty",
		},
		{
			name:  "Invalid alias",
			alias: "test",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			URLDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				if tc.alias != "" {
					URLDeleterMock.On("DeleteURL", context.Background(), tc.alias).
						Return(tc.mockError)
				}
			}

			log := sl.SetupLogger("dev")
			r := chi.NewRouter()
			r.Delete("/{alias}", Delete(log, URLDeleterMock))
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tc.alias), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if tc.alias == "" {
				require.Equal(t, http.StatusNotFound, rr.Code)
			} else {
				require.Equal(t, http.StatusOK, rr.Code)

				body := rr.Body.String()

				var resp response.Response

				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				require.Equal(t, tc.respError, resp.Error)
			}
		})
	}
}
