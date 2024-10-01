package httpkit

import (
	"net/http"

	"github.com/pawpawchat/core/pkg/response"
)

func BadRequest(w http.ResponseWriter, body map[string]any) {
	response.Json().BadRequest().Body(body).MustWrite(w)
}

func InternalError(w http.ResponseWriter) {
	response.Json().InternalError().Body(map[string]any{"error": "internal server error"}).MustWrite(w)
}

func Created(w http.ResponseWriter, body map[string]any) {
	response.Json().Created().Body(body).MustWrite(w)
}

func Ok(w http.ResponseWriter, body map[string]any) {
	response := response.Json().OK()
	if body != nil {
		response.Body(body)
	}
	response.MustWrite(w)
}
