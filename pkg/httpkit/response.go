package httpkit

import (
	"net/http"

	"github.com/pawpawchat/core/pkg/response"
)

func sendResponse(w http.ResponseWriter, code int, data any) {
	if data != nil {
		response.Json().Code(code).Body(data).MustWrite(w)
		return
	}
	response.Json().Code(code).MustWrite(w)
}

func SendWithCode(w http.ResponseWriter, code int, data any) {
	sendResponse(w, code, data)
}

func BadRequest(w http.ResponseWriter, data any) {
	sendResponse(w, http.StatusBadRequest, data)
}

func InternalError(w http.ResponseWriter, data any) {
	sendResponse(w, http.StatusInternalServerError, data)
}

func Created(w http.ResponseWriter, data any) {
	sendResponse(w, http.StatusCreated, data)
}

func Ok(w http.ResponseWriter, data any) {
	sendResponse(w, http.StatusOK, data)
}
