package redirect_test

import (
	"github.com/audetv/urlshortener/internal/http-server/handlers/redirect"
	"github.com/audetv/urlshortener/internal/http-server/handlers/redirect/mocks"
	"github.com/audetv/urlshortener/internal/lib/api"
	"github.com/audetv/urlshortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
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
			url:   "https://svodd.ru/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			//Check the final URL after redirection
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}