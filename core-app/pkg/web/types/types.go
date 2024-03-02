package types

import "net/http"

type Response[T any] struct {
	Success bool   `json:"status"`
	Error   string `json:"error"`
	Data    T      `json:"data,omitempty"`
}

func SuccessResponse(data any) (int, Response[any]) {
	return http.StatusOK, Response[any]{
		Success: true,
		Data:    data,
	}
}

func InternalErrorResponse(err error) (int, Response[string]) {
	return http.StatusInternalServerError, Response[string]{
		Success: false,
		Error:   err.Error(),
	}
}
