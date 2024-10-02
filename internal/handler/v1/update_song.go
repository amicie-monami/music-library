package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/pawpawchat/core/pkg/response"
)

type songDataUpdater interface {
	Tx(ctx context.Context, txActions func() error) error
	UpdateSong(ctx context.Context, song *model.Song) (count int64, err error)
	UpdateSongDetails(ctx context.Context, details *model.SongDetail) (count int64, err error)
}

// @Summary Изменение данных песни
// @Description Метод позволяет изменить данные песни, хранящиеся в библиотеке.
// @Router /songs/{id} [patch]
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор песни, данные которой необходимо изменить."
// @Param songInfo body dto.UpdateSongRequest true "Данные песни, которые необходимо изменить."
// @Success 200 {string} string "Данные были успешно обновлены, нет возвращаемого значения."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func UpdateSong(repo songDataUpdater) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songID, songIDVar, err := parsePathVarSongID(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": "invalid song id", "song_id": songIDVar}).MustWrite(w)
			return
		}

		//parse body
		song, songDetails, err := parseUpdateSongBody(songID, r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": err.Error()}).MustWrite(w)
			return
		}

		//transaction actions
		tx := func() error {
			if song != nil {
				count, err := repo.UpdateSong(r.Context(), song)
				if err != nil {
					slog.Error("failed to update a song data", "msg", err.Error())
					// needs refactoring: to get rid of the "magic" error
					response.Json().InternalError().Body(dto.Error{Message: "Internal server error"}).MustWrite(w)
					return err

				} else if count == 0 {
					slog.Info("song not found", "song_id", songID)
					response.Json().BadRequest().Body(body{"error": "song not found", "song_id": songID}).MustWrite(w)
					return err
				}
			}

			//sql query for song details
			if songDetails != nil {
				count, err := repo.UpdateSongDetails(r.Context(), songDetails)
				if err != nil {
					slog.Info("failed to update a song details", "msg", err.Error())
					// needs refactoring: to get rid of the "magic" error
					response.Json().InternalError().Body(dto.Error{Message: "Internal server error"}).MustWrite(w)
					return err

				} else if count == 0 {
					slog.Info("song details not found", "song_id", songID)
					response.Json().BadRequest().Body(body{"error": "song details not found", "song_id": songID}).MustWrite(w)
					return err
				}
			}
			return nil
		}

		if err := repo.Tx(r.Context(), tx); err != nil {
			return
		}

		response.Json().OK().MustWrite(w)
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
			SongID: songID,
			Text:   body.SongDetails.Text,
			Link:   body.SongDetails.Link,
		}
		if body.SongDetails.ReleaseDate != "" {
			releaseDate, err := time.Parse("01.01.2006", body.SongDetails.ReleaseDate)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse release_date field")
			}
			songDetails.ReleaseDate = &releaseDate
		}
	}

	return song, songDetails, nil
}
