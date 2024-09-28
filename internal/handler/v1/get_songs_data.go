package handler

import (
	"log/slog"
	"net/http"
	"strconv"
)

type songDataGetter interface {
}

func GetSongsData(repo songDataGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, err := GetIntParam("offset", r, true)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		limit, err := GetIntParam("limit", r, true)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		_ = offset
		_ = limit
	})
}

func GetIntParam(key string, r *http.Request, errlog bool) (int64, error) {
	if params := r.URL.Query(); params.Get(key) != "" {
		param, err := strconv.ParseInt(params.Get(key), 0, 10)
		if err != nil {
			if errlog {
				slog.Info("failed to parse param", "key", key, "err", err)
			}
			return 0, err
		}
		return param, nil
	}
	return 0, nil
}
