package handler

import (
	"encoding/json"
	"fmt"
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
		songId, err := httpkit.GetIntParam("id", r)
		if err != nil {
			slog.Info("failed to parse song id", "err", err)
			w.Write([]byte(err.Error()))
			return
		}

		offset, err := httpkit.GetIntParam("offset", r)
		if err != nil {
			slog.Info("failed to parse param 'offset'", "err", err)
			w.Write([]byte(err.Error()))
			return
		}

		limit, err := httpkit.GetIntParam("limit", r)
		if err != nil {
			slog.Info("failed to parse param 'limit'", "err", err)
			w.Write([]byte(err.Error()))
			return
		}

		songText, err := repo.GetSongText(songId)
		if err != nil {
			slog.Info("failed", "err", err)
			w.Write([]byte(err.Error()))
			return
		}

		if songText == nil {
			slog.Info("success", "err", err)
			w.Write([]byte("no data about song text"))
			return
		}

		couplets, err := textPagination(*songText, offset, limit)
		if err != nil {
			slog.Info("pagination failed", "err", err)
			w.Write([]byte(err.Error()))
			return
		}

		response, _ := json.Marshal(map[string]any{"couplets": couplets})
		slog.Info("success", "status_code", http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})
}

func textPagination(text string, offset int64, limit int64) ([]string, error) {
	if limit == 0 {
		limit = 1
	}

	couplets := strings.Split(text, "\n\n")
	if couplets[len(couplets)-1] == "" {
		couplets = couplets[:len(couplets)-1]
	}

	lenVerses := len(couplets)
	if lenVerses <= int(offset) {
		return nil, fmt.Errorf("out of bounds")
	}

	if int(offset+limit) >= lenVerses {
		return couplets[offset:], nil
	}

	return couplets[offset : offset+limit], nil
}
