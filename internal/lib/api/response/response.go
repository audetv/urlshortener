// Package response Здесь два привычных поля — Status и Error.
// Как и во многих других API-сервисах, эти поля могут присутствовать в ответе любого хэндлера.
// А раз так, то имеет смысл их вынести в общий пакет, он будет тут: internal/lib/api/response.
//
// Также я завел константы, которыми будем заполнять поле Status:
package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// Error Код можно сделать немного красивее, если вынести повторяющийся код
// формирования объекта ответа в общую функцию.
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func OK(msg string) Response {
	return Response{
		Status: StatusOK,
	}
}
