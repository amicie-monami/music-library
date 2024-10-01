package handler

import (
	"log/slog"
	"net/http"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songDetailsGetter interface {
	GetSongWithDetails(group string, title string) (*dto.SongWithDetails, error)
}

func GetSongDetails(repo songDetailsGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongDetailsQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": err.Error()})
			return
		}

		songWithDetails, err := repo.GetSongWithDetails(params["song"], params["title"])
		if err != nil {
			slog.Error(err.Error())
			httpkit.InternalError(w)
			return
		}

		httpkit.Ok(w, map[string]any{"song": songWithDetails})
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
