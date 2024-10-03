package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amicie-monami/music-library/internal/domain/mock"
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetSongs(t *testing.T) {
	testCases := []struct {
		Description string
		QueryParams string
		Code        int
	}{
		{
			Description: "Valid filter param",
			QueryParams: "filter=release_date=01.02.2022-08.02.2024,groups=ping+pong",
			Code:        http.StatusOK,
		},
		{
			Description: "Valid fields param",
			QueryParams: "fields=song_id",
			Code:        http.StatusOK,
		},
		{
			Description: "Valid offset param",
			QueryParams: "offset=12",
			Code:        http.StatusOK,
		},
		{
			Description: "Invalid limit param",
			QueryParams: "limit=12",
			Code:        http.StatusOK,
		},
		{
			Description: "Invalid filter param",
			QueryParams: "filter=catch",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid fields param",
			QueryParams: "fields=...",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid limit param",
			QueryParams: "limit=...",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid limit param",
			QueryParams: "limit=-1",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid offset param",
			QueryParams: "offset=...",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid offset param",
			QueryParams: "offset=-2",
			Code:        http.StatusBadRequest,
		},
	}

	getSongDetailsHandler := handler.GetSongs(&mock.SongRepo{})

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {

			urlWithQueryParams := fmt.Sprintf("/api/v1/songs?%s", tc.QueryParams)

			request := httptest.NewRequest("GET", urlWithQueryParams, nil)

			rr := httptest.NewRecorder()

			getSongDetailsHandler.ServeHTTP(rr, request)

			assert.Equal(t, tc.Code, rr.Code)
		})
	}
}
