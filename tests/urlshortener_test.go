package tests

import (
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/save"
	"github.com/audetv/urlshortener/internal/lib/api"
	"github.com/audetv/urlshortener/internal/lib/random"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

const host = "localhost:8082"

// TestURLShortener_HappyPath tests the happy path scenario of the URLShortener function.
//
// It creates a URL using the provided scheme and host, and then sends a POST request to the "/url" path.
// The request body is formed using a randomly generated URL and alias. The function then authenticates using basic authentication.
// The response is expected to have a status code of 200, and the response body is expected to be in JSON format.
// The response body is also expected to contain the key "alias".
// Чтобы выполнить тест, нужно сначала запустить сервис, затем уже — тест.
func TestURLShortener_HappyPath(t *testing.T) {
	// Универсальный способ создать URL
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// Создаем экземпляр клиент httpexpect для тестирования
	e := httpexpect.Default(t, u.String())

	e.POST("/url"). // Отправляем POST-запрос, путь - "/url"
			WithJSON(save.Request{ // Формируем тело запроса
			URL:   gofakeit.URL(),             // Генерируем случайный URL
			Alias: random.NewRandomString(10), // Генерируем случайный алиас
		}).
		WithBasicAuth("user", "password"). // Авторизуемся
		Expect().                          // Далее перечисляем наши ожидания от ответа
		Status(200).                       // Ожидаем код статуса 200
		JSON().
		Object().            // Ожидаем, что в теле ответа будет JSON
		ContainsKey("alias") // Проверяем, что в теле есть ключ "alias"
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		alias string
		url   string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("user", "password").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").IsEqual(tc.error)
				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}
	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
