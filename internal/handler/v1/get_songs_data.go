package handler

import (
	"net/http"
	"strconv"
)

type songDataGetter interface {
}

func GetSongsData(repo songDataGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, err := GetIntParam("offset", r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		limit, err := GetIntParam("limit", r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		_ = offset
		_ = limit

	})
}

func GetIntParam(key string, r *http.Request) (int64, error) {
	if params := r.URL.Query(); params.Get(key) != "" {
		param, err := strconv.ParseInt(params.Get(key), 0, 10)
		if err != nil {
			return 0, err
		}
		return param, nil
	}
	return 0, nil
}
