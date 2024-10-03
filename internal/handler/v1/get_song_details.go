package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
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
// @Success 201 {object} dto.GetSongDetailsResponse "Объект, описывающий основную и дополнительную информацию о песне."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func GetSongDetails(repo songDetailsGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongDetailsQueryParams(r)
		if err != nil {
			sendError(w, err)
			return
		}

		songWithDetails, err := repo.GetSongWithDetails(r.Context(), params["group"], params["song"])
		if err != nil {
			sendError(w, err)
			return
		}

		slog.Info("song has been found", "id", songWithDetails.ID)
		responseBody := dto.GetSongDetailsResponse{Song: songWithDetails}
		httpkit.Ok(w, responseBody)
	})
}

func parseGetSongDetailsQueryParams(r *http.Request) (map[string]string, error) {
	params := make(map[string]string)
	var err error

	params["group"], err = httpkit.GetStrRequiredParam("group", r)
	if err != nil {
		return nil, dto.NewError(400, "group param is required", "parseGetSongDetailsQueryParams", nil, nil)
	}

	params["song"], err = httpkit.GetStrRequiredParam("song", r)
	if err != nil {
		return nil, dto.NewError(400, "song param is required", "parseGetSongDetailsQueryParams", nil, nil)
	}

	return params, nil
}
