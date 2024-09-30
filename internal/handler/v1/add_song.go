package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/model"
)

type SongAdder interface {
	Create(song *model.Song) error
}

// AddSong ...
func AddSong(repo SongAdder) http.Handler {

	type AddSongRequest struct {
		Group string `json:"group"`
		Title string `json:"title"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//parse the request body
		var data AddSongRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			slog.Info("failed to parse body", "msg", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//form song model
		song := &model.Song{Group: data.Group, Name: data.Title}

		//add song to database
		if err := repo.Create(song); err != nil {
			slog.Info("failed to add the song to a db", "msg", err.Error())
			// refactor -- checks error source, maybe it's http.ServerInternalError
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//marshal a created song
		response, err := json.Marshal(song)
		if err != nil {
			slog.Info("failed to marshal song model to json", "msg", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//response
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
		slog.Info("success", "status_code", http.StatusCreated)
	})
}
