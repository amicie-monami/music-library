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
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songDataUpdater interface {
	Tx(ctx context.Context, txActions func() error) error
	UpdateSong(ctx context.Context, song *model.Song) error
	UpdateSongDetails(ctx context.Context, details *model.SongDetail) error
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
		songID, err := parsePathVarSongID(r)
		if err != nil {
			sendError(w, err)
			return
		}

		//parse body
		song, songDetails, err := parseUpdateSongBody(songID, r)
		if err != nil {
			sendError(w, err)
			return
		}

		//transaction actions
		tx := func() error {
			if song != nil && repo.UpdateSong(r.Context(), song) != nil {
				return err
			}

			if songDetails != nil && repo.UpdateSongDetails(r.Context(), songDetails) != nil {
				return err
			}
			return nil
		}

		if err := repo.Tx(r.Context(), tx); err != nil {
			sendError(w, err)
			return
		}

		slog.Info("song has been successfully updated", "id", songID)
		httpkit.Ok(w, nil)
	})
}

func parseUpdateSongBody(songID int64, r *http.Request) (*model.Song, *model.SongDetail, error) {
	var (
		requestBody dto.UpdateSongRequest
		song        *model.Song
		songDetails *model.SongDetail
	)

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, nil, err
	}

	if requestBody.Group == "" && requestBody.Song == "" && requestBody.ReleaseDate == "" && requestBody.Link == "" && requestBody.Text == "" {
		return nil, nil, dto.NewError(400, "missing the data for updates", "parseUpdateSongBody", nil, nil)
	}

	song = &model.Song{ID: songID, Group: requestBody.Group, Name: requestBody.Song}

	songDetails = &model.SongDetail{SongID: songID, Text: &requestBody.Text, Link: &requestBody.Link}

	//if release_date has been sent
	if requestBody.ReleaseDate != "" {

		//try to parse release_date
		releaseDate, err := time.Parse("01.01.2006", requestBody.ReleaseDate)
		if err != nil {
			details := fmt.Sprintf("release_date=%s", requestBody.ReleaseDate)
			return nil, nil, &dto.Error{Code: 400, Message: "failed to parse release_date field", Details: details}
		}

		songDetails.ReleaseDate = &releaseDate
	}

	return song, songDetails, nil
}
