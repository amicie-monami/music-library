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

func TestGetSongWithDetails(t *testing.T) {
	testCases := []struct {
		Description string
		QueryParams string
		Code        int
	}{
		{
			Description: "Valid query params",
			QueryParams: fmt.Sprintf("group=%s&song=%s", mock.ValidGroupName, mock.ValidSongName),
			Code:        http.StatusOK,
		},
		{
			Description: "Missing the required query param group",
			QueryParams: "group=dummy",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Missing the required query param song",
			QueryParams: "song=dummy",
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Both required params are missing",
			QueryParams: "",
			Code:        http.StatusBadRequest,
		},
	}

	getSongDetailsHandler := handler.GetSongDetails(&mock.SongRepo{})

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/songs?%s", tc.QueryParams)
			// fmt.Println(url)
			request := httptest.NewRequest("GET", url, nil)

			requestRecorder := httptest.NewRecorder()

			getSongDetailsHandler.ServeHTTP(requestRecorder, request)

			assert.Equal(t, tc.Code, requestRecorder.Code)
		})
	}
}
