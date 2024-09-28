package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SongDeletter interface {
	Delete(id int64) error
}

// DeleteSobg ...
func DeleteSong(repo SongDeletter) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		songIdStr := mux.Vars(r)["id"]
		songId, err := strconv.ParseInt(songIdStr, 0, 10)
		if err != nil {
			slog.Info("failed to parse song id", "value", songIdStr)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//add song to database
		if err := repo.Delete(songId); err != nil {
			slog.Info("failed to delete a song", "msg", err.Error())
			// refactor -- checks error source, maybe it's http.ServerInternalError
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//response
		w.WriteHeader(http.StatusOK)
		slog.Info("success", "status_code", http.StatusOK)
	})
}
