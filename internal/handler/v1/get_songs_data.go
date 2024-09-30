package handler

import (
	"encoding/json"
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

// GetSongsData ...
func GetSongsData(usecase songDataGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongsDataQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		data, err := usecase.GetSongs(params)
		if err != nil {
			slog.Info(err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		response, _ := json.Marshal(data)
		slog.Info("success", "status_code", http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})
}

// parseGetSongsDataQueryParams ...
func parseGetSongsDataQueryParams(r *http.Request) (map[string]any, error) {
	filterMap, err := parseFilterParams(r)
	if err != nil {
		return nil, err
	}

	paginMap, err := parsePaginationParams(r)
	if err != nil {
		return nil, err
	}

	fields := httpkit.GetStrParam("fields", r)
	return map[string]any{"filter": filterMap, "pagination": paginMap, "fields": fields}, nil
}

// parsePaginationParams ...
func parsePaginationParams(r *http.Request) (map[string]int64, error) {
	offset, err := httpkit.GetIntParam("offset", r)
	if err != nil {
		return nil, err
	}

	offset, err = checkOffsetParam(offset)
	if err != nil {
		return nil, err
	}

	limit, err := httpkit.GetIntParam("limit", r)
	if err != nil {
		return nil, err
	}

	limit, err = checkLimitParam(limit)
	if err != nil {
		return nil, err
	}

	return map[string]int64{"offset": offset, "limit": limit}, nil
}

func checkLimitParam(limit int64) (int64, error) {
	if limit == 0 {
		return 10, nil
	}

	if limit > 1000 {
		return 1000, nil
	}

	if limit < 1 {
		return 0, fmt.Errorf("limit=%d, but param must be >= 1", limit)
	}

	return limit, nil
}

func checkOffsetParam(offset int64) (int64, error) {
	if offset < 0 {
		return 0, fmt.Errorf("offset=%d, but param must be >= 0", offset)
	}

	return offset, nil
}

// parseFilterParams
func parseFilterParams(r *http.Request) (map[string]any, error) {
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
