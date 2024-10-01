package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type SongAdder interface {
	Create(song *model.Song) error
}

func AddSong(repo SongAdder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		song, err := parseAddSongBody(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": err.Error()})
			return
		}

		if err := repo.Create(song); err != nil {
			slog.Error(err.Error())
			httpkit.InternalError(w)
			return
		}

		httpkit.Created(w, map[string]any{"added_song": song})
		slog.Info("success", "status_code", http.StatusCreated)
	})
}

func parseAddSongBody(r *http.Request) (*model.Song, error) {
	var data dto.AddSongRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.Group == "" {
		return nil, fmt.Errorf("field group is required")
	}

	if data.Song == "" {
		return nil, fmt.Errorf("field song is required")
	}

	return &model.Song{Group: data.Group, Name: data.Song}, nil
}
