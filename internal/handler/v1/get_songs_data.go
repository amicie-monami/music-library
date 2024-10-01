package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songDataGetter interface {
	GetSongs(aggregation map[string]any) ([]dto.SongWithDetails, error)
}

func GetSongsData(repo songDataGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongsDataQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			httpkit.BadRequest(w, map[string]any{"error": err.Error()})
			return
		}

		data, err := repo.GetSongs(params)
		if err != nil {
			slog.Info(err.Error())
			httpkit.InternalError(w)
			return
		}

		slog.Info("200")
		httpkit.Ok(w, map[string]any{"songs": data})
	})
}

func parseGetSongsDataQueryParams(r *http.Request) (map[string]any, error) {
	filterMap, err := parseGetSongsDataFilterParams(r)
	if err != nil {
		return nil, err
	}

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		return nil, err
	}

	fields := httpkit.GetStrParam("fields", r)
	return map[string]any{
		"filter": filterMap,
		"limit":  limit,
		"offset": offset,
		"fields": fields,
	}, nil
}

func parseGetSongsDataFilterParams(r *http.Request) (map[string]any, error) {
	filter := httpkit.GetStrParam("filter", r)
	if filter == "" {
		return nil, nil
	}

	filterMap := make(map[string]any)
	params := strings.Split(filter, ",")

	for idx := range params {
		param := strings.Split(params[idx], "=")
		if len(param) != 2 {
			return nil, fmt.Errorf("invalid param %s", params[idx])
		}
		filterMap[param[0]] = param[1]
	}

	return filterMap, nil
}
