package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songDataUpdater interface {
	Tx(txActions func() error) error
	UpdateSong(song *model.Song) (count int64, err error)
	UpdateSongDetails(details *model.SongDetail) (count int64, err error)
}

func UpdateSong(repo songDataUpdater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songID, songIDVar, err := parsePathVarSongID(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": "invalid song id", "song_id": songIDVar})
			return
		}

		//parse body
		song, songDetails, err := parseUpdateSongBody(songID, r)
		if err != nil {
			slog.Info("404" + err.Error())
			httpkit.BadRequest(w, map[string]any{"error": err.Error()})
			return
		}

		//refactor: wrap in transaction
		//sql query for song
		count, err := repo.UpdateSong(song)
		if err != nil {
			slog.Error("404 failed to update a song data", "msg", err.Error())
			httpkit.InternalError(w)
			return
		} else if count == 0 {
			slog.Info("404 song not found", "song_id", songID)
			httpkit.BadRequest(w, map[string]any{"error": "song not found", "song_id": songID})
			return
		}

		//sql query for song details
		count, err = repo.UpdateSongDetails(songDetails)
		if err != nil {
			slog.Info("500 failed to update a song details", "msg", err.Error())
			httpkit.InternalError(w)
			return
		} else if count == 0 {
			slog.Info("404 song details not found", "song_id", songID)
			httpkit.BadRequest(w, map[string]any{"error": "song not found", "song_id": songID})
			return
		}

		slog.Info("200")
	})
}

func parseUpdateSongBody(songID int64, r *http.Request) (*model.Song, *model.SongDetail, error) {
	var (
		body        dto.UpdateSongRequest
		song        *model.Song
		songDetails *model.SongDetail
	)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, nil, err
	}

	if body.Song != nil {
		song = &model.Song{
			ID:    songID,
			Group: body.Song.Group,
			Name:  body.Song.Name,
		}
	}

	if body.SongDetails != nil {
		songDetails = &model.SongDetail{
			SongID:      songID,
			ReleaseDate: body.SongDetails.ReleaseDate,
			Text:        body.SongDetails.Text,
			Link:        body.SongDetails.Link,
		}
	}

	return song, songDetails, nil
}
