package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/save"
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/save/mocks"
	"github.com/audetv/urlshortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string // Имя теста
		alias     string // Отправляемый alias
		url       string // Отправляемый URL
		respError string // Какую ошибку мы должны получить?
		mockError error  // Ошибку, которую вернет mock
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://svodd.ru",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://svodd.ru",
		},
		// Другие кейсы
	}

	for _, tc := range cases {
		tc := tc // Create a new variable to avoid closure issues
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Создаём объект мока стораджа
			urlSaverMock := mocks.NewURLSaver(t)

			// Если ожидается успешный ответ, значит к моку точно будет вызов
			// Либо даже если в ответе ожидаем ошибку
			// но мок должен ответить с ошибкой, к нему тоже будет запрос:
			if tc.respError == "" || tc.mockError != nil {
				// Сообщеам моку, какой к нему будет запрос, и что надо вернуть
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once() // Запрос будет ровно один
			}

			// Создаём наш хэндлер
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			// Формируем тело запроса
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			// Создаём объект запроса
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Создаём ResponseRecorder для записи ответа хэндлера
			rr := httptest.NewRecorder()
			// Обрабатываем запрос,записывая ответ в рекордер
			handler.ServeHTTP(rr, req)

			// Проверяем, что статус ответа корректный
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем наличие требуемой ошибки в ответе
			require.Equal(t, tc.respError, resp.Error)

			// Другие проверки
		})
	}
}
