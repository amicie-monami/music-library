package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
	"github.com/pawpawchat/core/pkg/response"
)

type songTextGetter interface {
	GetSongText(id int64) (*string, error)
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
		songID, songIDVar, err := parsePathVarSongID(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": "invalid song id", "song_id": songIDVar}).MustWrite(w)
			return
		}

		limit, offset, err := parseGetSongTextQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": err.Error()}).MustWrite(w)
			return
		}

		songText, err := repo.GetSongText(songID)
		if err != nil {

			if err == sql.ErrNoRows {
				slog.Info("400 No data")
				response.Json().
					BadRequest().
					Body(dto.Error{
						Message: "There is no information about the lyrics of a song with this identifier",
						Details: fmt.Sprintf("song_id=%s", songIDVar),
					}).
					MustWrite(w)
				return
			}

			slog.Error("500: " + err.Error())
			// needs refactoring: to get rid of the "magic" error
			response.Json().InternalError().Body(dto.Error{Message: "Internal server error"})
			return
		}

		couplets, err := coupletsPagination(songText, limit, offset)
		if err != nil {
			response.Json().BadRequest().Body(body{"error": err.Error()}).MustWrite(w)
			return
		}

		slog.Info("200")
		response.Json().OK().Body(dto.GetSongTextResponse{Couplets: couplets, SongID: songID}).MustWrite(w)
	})
}

func parseGetSongTextQueryParams(r *http.Request) (int64, int64, error) {
	limit, err := httpkit.GetIntParam("limit", r)
	if err != nil {
		return 0, 0, err
	}
	offset, err := httpkit.GetIntParam("offset", r)
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
