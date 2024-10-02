package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
	"github.com/pawpawchat/core/pkg/response"
)

type songDetailsGetter interface {
	GetSongWithDetails(ctx context.Context, group string, title string) (*dto.SongWithDetails, error)
}

// @Summary Информация о песне
// @Description Метод возвращает полную информацию о песне.
// @Router /info [get]
// @Tags Songs
// @Produce json
// @Param group query string true "Название группы"
// @Param song query string  true "Название песни"
// @Success 201 {object} dto.SongWithDetails "Объект, описывающий основную и дополнительную информацию о песне."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func GetSongDetails(repo songDetailsGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongDetailsQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": err.Error()}).MustWrite(w)
			return
		}

		songWithDetails, err := repo.GetSongWithDetails(r.Context(), params["song"], params["title"])
		if err != nil {
			slog.Error(err.Error())
			// needs refactoring: to get rid of the "magic" error
			response.Json().InternalError().Body(dto.Error{Message: "Internal server error"})
			return
		}

		response.Json().OK().Body(map[string]any{"song": songWithDetails}).MustWrite(w)
		slog.Info("200")
	})
}

func parseGetSongDetailsQueryParams(r *http.Request) (map[string]string, error) {
	params := make(map[string]string)
	var err error

	params["group"], err = httpkit.GetStrRequiredParam("group", r)
	if err != nil {
		return nil, err
	}

	params["song"], err = httpkit.GetStrRequiredParam("song", r)
	if err != nil {
		return nil, err
	}

	return params, nil
}
