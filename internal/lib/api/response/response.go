// Package response Здесь два привычных поля — Status и Error.
// Как и во многих других API-сервисах, эти поля могут присутствовать в ответе любого хэндлера.
// А раз так, то имеет смысл их вынести в общий пакет, он будет тут: internal/lib/api/response.
//
// Также я завел константы, которыми будем заполнять поле Status:
package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

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

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// ValidationError Для формирование более ясного ответа добавляем в пакет response такую функцию:
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
