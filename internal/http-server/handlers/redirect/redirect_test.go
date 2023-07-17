package redirect_test

import (
	"context"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/redirect"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/redirect/mocks"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/api"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				if tc.alias != "" { // Проверка, что alias не пустой перед вызовом GetURL
					urlGetterMock.On("GetURL", context.Background(), tc.alias).
						Return(tc.url, tc.mockError).Once()
				}
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			if tc.alias == "" { // Проверка, что alias пустой
				_, err := api.GetRedirect(ts.URL + "/" + tc.alias)

				assert.Error(t, err)
				return
			}

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			assert.NoError(t, err)

			// Check the final URL after redirection.
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
