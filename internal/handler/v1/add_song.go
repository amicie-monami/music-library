package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/pawpawchat/core/pkg/response"
)

type body map[string]any

type SongAdder interface {
	Create(song *model.Song) error
}

// @Summary Добавление новой песни
// @Description Метод добавляет в библиотеку основную информацию о песне.
// @Router /songs [post]
// @Tags Songs
// @Accept json
// @Produce json
// @Param group body dto.AddSongRequest true "Параметры песни, информацию о которой необходимо добавить в библиотеку."
// @Success 201 {object} model.Song "Объект, описывающий добавленную песню."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func AddSong(repo SongAdder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		song, err := parseAddSongBody(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": err.Error()}).MustWrite(w)
			return
		}

		if err := repo.Create(song); err != nil {
			slog.Error(err.Error())
			// needs refactoring: to get rid of the "magic" error
			response.Json().InternalError().Body(body{"error": "internal server error"}).MustWrite(w)
			return
		}

		slog.Info("201")
		response.Json().InternalError().Body(body{"added_song": song}).MustWrite(w)
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
