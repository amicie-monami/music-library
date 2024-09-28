package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/model"
)

type SongAdder interface {
	Create(song *model.Song) error
}

func AddSong(repo SongAdder) http.Handler {

	type AddSongRequest struct {
		Group string `json:"group"`
		Title string `json:"title"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//parse the request body
		var data AddSongRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			slog.Debug("failed to parse body", "msg", err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		//form song model
		song := &model.Song{
			Group: data.Group,
			Title: data.Title,
		}

		//add song to database
		if err := repo.Create(song); err != nil {
			slog.Debug("failed to add the song to a db", "msg", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		//marshal created song to json
		response, err := json.Marshal(song)
		if err != nil {
			slog.Debug("failed to marshal song model to json", "msg", err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		//response
		w.WriteHeader(201)
		w.Write(response)
		slog.Debug("success", "status_code", 201)
	})
}
