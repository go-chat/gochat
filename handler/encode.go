package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chat/gochat/apperror"
)

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

func encodeErrorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	r := &Response{
		Success: false,
		Error:   err.Error(),
	}
	encodeResponse(w, r, httpStatusCode)
}

func encodeAppErrorResponse(w http.ResponseWriter, appError *apperror.AppError) {
	r := &Response{
		Success: false,
		Error:   appError.Message,
	}
	encodeResponse(w, r, appError.StatusCode)
}

func encodeSuccessResponse(w http.ResponseWriter, data interface{}) {
	r := &Response{
		Success: true,
		Error:   "",
		Data:    data,
	}

	encodeResponse(w, r, http.StatusOK)
}

func encodeResponse(w http.ResponseWriter, data interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
}
