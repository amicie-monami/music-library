package handler

import (
	"log/slog"
	"net/http"

	"github.com/pawpawchat/core/pkg/response"
)

type SongDeletter interface {
	Delete(id int64) (int64, error)
}

// @Summary Удаление песни
// @Description Метод удаляет всю информацию о песне по переданному идектификатору.
// @Router /songs/{id} [delete]
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор песни, информацию о которой необходимо удалить."
// @Success 200 {string} string "Информация успешно удалена, нет данных в теле ответа."
// @Failure 404 {object} dto.Error "Неккоректные значения параметров запроса."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func DeleteSong(repo SongDeletter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songID, songIDVar, err := parsePathVarSongID(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": "invalid song id", "song_id": songIDVar}).MustWrite(w)
			return
		}

		count, err := repo.Delete(songID)
		if err != nil {
			slog.Error(err.Error())
			// needs refactoring: to get rid of the "magic" error
			response.Json().InternalError().Body(body{"error": "internal server error"}).MustWrite(w)
			return
		}

		if count == 0 {
			slog.Info("not found song id=" + songIDVar)
			response.Json().BadRequest().Body(body{"error": "song not found", "song_id": songID}).MustWrite(w)
			return
		}

		response.Json().OK().MustWrite(w)
		slog.Info("200")
	})
}
