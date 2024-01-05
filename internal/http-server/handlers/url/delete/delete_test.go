package delete_test

import (
	"encoding/json"
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/delete"
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/delete/mocks"
	"github.com/audetv/urlshortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO Сделать рабочий тест, не работает, не разбирает параметр alias
func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)
			r := chi.NewRouter()
			r.Delete("/url/{alias}", handler)

			req := httptest.NewRequest("DELETE", "/url/"+tc.alias, nil)
			// Создаём ResponseRecorder для записи ответа хэндлера
			rr := httptest.NewRecorder()
			// Обрабатываем запрос,записывая ответ в рекордер
			handler.ServeHTTP(rr, req)

			// Проверяем, что статус ответа корректный
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d but got %d", http.StatusOK, rr.Code)
			}

			body := rr.Body.String()

			log.Println(body)

			var resp delete.Response

			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем наличие требуемой ошибки в ответе
			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
