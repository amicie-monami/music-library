package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type SongAdder interface {
	Create(ctx context.Context, song *model.Song) error
}

// @Summary Добавление новой песни
// @Description Метод добавляет в библиотеку основную информацию о песне.
// @Router /songs [post]
// @Tags Songs
// @Accept json
// @Produce json
// @Param group body dto.AddSongRequest true "Параметры песни, информацию о которой необходимо добавить в библиотеку."
// @Success 201 {object} dto.AddSongResponse "Объект, описывающий добавленную песню."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func AddSong(repo SongAdder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		song, err := parseAddSongBody(r)
		if err != nil {
			sendError(w, err)
			return
		}

		if err := repo.Create(r.Context(), song); err != nil {
			sendError(w, err)
			return
		}

		slog.Info("song has been added", "id", song.ID, "group", song.Group, "song", song.Name)
		responseBody := dto.AddSongResponse{Song: &dto.Song{ID: song.ID, Group: song.Group, Name: song.Name}}
		httpkit.Created(w, responseBody)
	})
}

func parseAddSongBody(r *http.Request) (*model.Song, error) {
	var data dto.AddSongRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, dto.NewError(400, "failed to parse song data", "parseAddSongBody", err.Error(), nil)
	}

	if data.Group == "" {
		return nil, dto.NewError(400, "incorrect song data", "parseAddSongBody", "field group is required", nil)
	}

	if data.Song == "" {
		return nil, dto.NewError(400, "incorrect song data", "parseAddSongBody", "field song is required", nil)
	}

	return &model.Song{Group: data.Group, Name: data.Song}, nil
}

func sendError(w http.ResponseWriter, err error) {
	dtoErr, ok := err.(*dto.Error)
	if !ok {
		// unkonwn error - log level error
		slog.Error("unkown error", "type", reflect.TypeOf(err), "err", err.Error())
		httpkit.InternalError(w, &dto.Error{Code: 500, Message: "internal server error"})
		return
	}

	switch dtoErr.Code {
	case 400: // bad request - log level info
		slog.Info(err.Error())
		httpkit.BadRequest(w, err)

	case 500: // internal server - log level error
		slog.Error(err.Error())
		httpkit.InternalError(w, err)

	default: // unknown code - log level error
		slog.Error("unkown error", "code", dtoErr, "err", err.Error())
		httpkit.InternalError(w, &dto.Error{Code: 500, Message: "internal server error"})
		return
	}
}
