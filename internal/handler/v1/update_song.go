package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amicie-monami/music-library/internal/model"
	"github.com/gorilla/mux"
)

type songDataUpdater interface {
	Tx(txActions func() error) error
	UpdateSong(song *model.Song) error
	UpdateSongDetails(details *model.SongDetail) error
}

func UpdateSong(repo songDataUpdater) http.Handler {

	type UpdateSongDTO struct {
		Group string `json:"group,omitempty"`
		Title string `json:"title,omitempty"`
	}

	type SongDetailsDTO struct {
		Text        *string `json:"text,omitempty"`
		Link        *string `json:"link,omitempty"`
		ReleaseDate *string `json:"release_date,omitempty"`
	}

	type UpdateSongRequest struct {
		Song        *UpdateSongDTO  `json:"song"`
		SongDetails *SongDetailsDTO `json:"song_details"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songIdStr := mux.Vars(r)["id"]
		songId, err := strconv.ParseInt(songIdStr, 0, 10)
		if err != nil {
			slog.Info("invalid song id", "value", songIdStr)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//parse the request body
		var data UpdateSongRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			slog.Info("invalid request body", "msg", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if data.Song == nil && data.SongDetails == nil {
			slog.Info("missing data for updates")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing data for updates"))
			return
		}

		var song *model.Song
		var songDetails *model.SongDetail
		if data.Song != nil {
			song = &model.Song{ID: songId, Group: data.Song.Group, Title: data.Song.Title}
		}

		if data.SongDetails != nil {
			songDetails = &model.SongDetail{
				SongID:      songId,
				ReleaseDate: data.SongDetails.ReleaseDate,
				Text:        data.SongDetails.Text,
				Link:        data.SongDetails.Link,
			}
		}

		err = repo.Tx(func() error {
			if err := repo.UpdateSong(song); err != nil {
				slog.Info("failed to update a song data", "msg", err.Error())
				return err
			}

			if err := repo.UpdateSongDetails(songDetails); err != nil {
				slog.Info("failed to update a song details", "msg", err.Error())
				return err
			}
			return nil
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		//response
		w.WriteHeader(http.StatusOK)
		slog.Info("success", "status_code", http.StatusOK)
	})
}
