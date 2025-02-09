package response

import (
	"fmt"
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

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationError) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "requered":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
