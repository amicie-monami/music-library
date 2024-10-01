package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/amicie-monami/music-library/pkg/httpkit"
	"github.com/gorilla/mux"
)

func parsePathVarSongID(r *http.Request) (int64, string, error) {
	songIDVar := mux.Vars(r)["id"]
	songID, err := strconv.ParseInt(songIDVar, 0, 10)
	if err != nil {
		return 0, songIDVar, fmt.Errorf("invalid song id: %s", songIDVar)
	}
	return songID, songIDVar, nil
}

func parsePaginationParams(r *http.Request) (int64, int64, error) {
	offset, err := httpkit.GetIntParam("offset", r)
	if err != nil {
		return 0, 0, err
	}
	offset, err = proccessOffsetParam(offset)
	if err != nil {
		return 0, 0, err
	}

	limit, err := httpkit.GetIntParam("limit", r)
	if err != nil {
		return 0, 0, err
	}
	limit, err = proccessLimitParam(limit)
	if err != nil {
		return 0, 0, err
	}

	return limit, offset, nil
}

func proccessLimitParam(limit int64) (int64, error) {
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

func proccessOffsetParam(offset int64) (int64, error) {
	if offset < 0 {
		return 0, fmt.Errorf("offset=%d, but param must be >= 0", offset)
	}

	return offset, nil
}
