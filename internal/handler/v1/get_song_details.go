package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/model"
)

type songDetailsGetter interface {
	GetSongDetails(group string, title string) (*model.SongDetail, error)
}

func GetSongDetails(repo songDetailsGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		if group == "" {
			slog.Info("missing required parameter `group`")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing required parameter `group`"))
			return
		}

		title := r.URL.Query().Get("song")
		if title == "" {
			slog.Info("missing required parameter `name`")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing required parameter `song`"))
			return
		}

		details, err := repo.GetSongDetails(group, title)
		if err != nil {
			slog.Info("failed", "msg", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		response, err := json.Marshal(details)
		if err != nil {
			slog.Info("failed to marshal song details", "msg", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
		slog.Info("success", "status_code", http.StatusOK)
	})
}
