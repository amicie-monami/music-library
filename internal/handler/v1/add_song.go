package handler

import (
	"bytes"
	"net/http"
)

func AddSong() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bytes.NewBufferString("string").Bytes())
	})
}
