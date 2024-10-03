package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songTextGetter interface {
	GetSongText(ctx context.Context, id int64) (*string, error)
}

// @Summary Получение текста песни с пагинацией по куплетам
// @Router /songs/{id}/lyrics [get]
// @Description Метод возвращает текст песни в куплетах. Если не заданы параметры пагинации, возвращаются все куплеты.
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор песни, текст которой необходимо получить."
// @Param limit query string false "Количество куплетов, которое необходимо верунть."
// @Param offset query string false "Смещение, необходимое для выборки определенного подмножества куплетов."
// @Success 200 {object} dto.GetSongTextResponse "Текст песни"
// @Failure 400 {object} dto.Error "Неверный запрос, некорректые значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func GetSongText(repo songTextGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songID, err := parsePathVarSongID(r)
		if err != nil {
			sendError(w, err)
			return
		}

		limit, offset, err := parseGetSongTextQueryParams(r)
		if err != nil {
			sendError(w, err)
			return
		}

		songText, err := repo.GetSongText(r.Context(), songID)
		if err != nil {
			sendError(w, err)
			return
		}

		couplets, err := coupletsPagination(songText, limit, offset)
		if err != nil {
			sendError(w, err)
			return
		}

		slog.Info("song lyrics have been found", "song_id", songID)
		responseBody := dto.GetSongTextResponse{Couplets: couplets, SongID: songID}
		httpkit.Ok(w, responseBody)
	})
}

func parseGetSongTextQueryParams(r *http.Request) (int64, int64, error) {
	limit, err := parseLimitParam(r)
	if err != nil {
		return 0, 0, err
	}

	offset, err := parseOffsetParam(r)
	if err != nil {
		return 0, 0, err
	}

	return limit, offset, nil
}

func coupletsPagination(text *string, limit int64, offset int64) ([]string, error) {
	if text == nil {
		return nil, nil
	}

	couplets := strings.Split(strings.ReplaceAll(*text, `\n`, "\n"), "\n\n")
	if limit == 0 && offset == 0 {
		return couplets, nil
	}

	if limit == 0 {
		limit = 1000
	}

	if int(offset) >= len(couplets) {
		return nil, nil
	}

	if int(offset+limit) >= len(couplets) {
		return couplets[offset:], nil
	}

	return couplets[offset : offset+limit], nil
}
