package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songTextGetter interface {
	GetSongText(id int64) (*string, error)
}

func GetSongText(repo songTextGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songID, songIDVar, err := parsePathVarSongID(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": "invalid song id", "song_id": songIDVar})
			return
		}

		limit, offset, err := parseGetSongTextQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": err.Error()})
			return
		}

		songText, err := repo.GetSongText(songID)
		if err != nil {
			slog.Error(err.Error())
			httpkit.InternalError(w)
			return
		}

		couplets, err := coupletsPagination(songText, limit, offset)
		if err != nil {
			httpkit.BadRequest(w, map[string]any{"error": err.Error()})
			return
		}

		httpkit.Ok(w, map[string]any{"song_id": songID, "couplets": couplets})
		slog.Info("200")
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
