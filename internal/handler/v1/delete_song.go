package handler

import (
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type SongDeletter interface {
	Delete(id int64) (int64, error)
}

func DeleteSong(repo SongDeletter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songID, songIDVar, err := parsePathVarSongID(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": "invalid song id", "song_id": songIDVar})
			return
		}

		count, err := repo.Delete(songID)
		if err != nil {
			slog.Error(err.Error())
			httpkit.InternalError(w)
			return
		}

		if count == 0 {
			slog.Info("404 not found song id=" + songIDVar)
			httpkit.BadRequest(w, map[string]any{"error": "song not found", "song_id": songID})
			return
		}

		slog.Info("201")
	})
}
