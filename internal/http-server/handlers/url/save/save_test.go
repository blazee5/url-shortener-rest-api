package save

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/blazee5/url-shortener-rest-api"
	models "github.com/blazee5/url-shortener-rest-api"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/url/save/mocks"
	"github.com/blazee5/url-shortener-rest-api/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/require"
	"net/http"
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
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				if tc.url != "" {
					urlSaverMock.On("SaveURL", context.Background(), &models.ShortUrl{
						ID:  tc.alias,
						URL: tc.url,
					}).
						Return(tc.mockError).Once()
				}
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
