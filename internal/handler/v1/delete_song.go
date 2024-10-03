package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
	"github.com/gorilla/mux"
)

type SongDeletter interface {
	Delete(ctx context.Context, id int64) error
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
		songID, err := parsePathVarSongID(r)
		if err != nil {
			sendError(w, err)
			return
		}

		if err := repo.Delete(r.Context(), songID); err != nil {
			sendError(w, err)
			return
		}

		slog.Info("song has been deleted", "id", songID)
		httpkit.Ok(w, nil)
	})
}

func parsePathVarSongID(r *http.Request) (int64, error) {
	songID, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 10)

	if err != nil {
		details := fmt.Sprintf("id=%s", mux.Vars(r)["id"])
		return 0, dto.NewError(400, "invalid song id in url", "parsePathVarSongID", details, nil)
	}

	if songID <= 0 {
		details := fmt.Sprintf("id=%s but must me > 0", mux.Vars(r)["id"])
		return 0, dto.NewError(400, "invalid song id in url", "parsePathVarSongID", details, nil)
	}

	return songID, nil
}
