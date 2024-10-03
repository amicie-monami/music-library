package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amicie-monami/music-library/internal/domain/mock"
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetSongText(t *testing.T) {
	testCases := []struct {
		Description string
		SongID      int64
		Code        int
	}{
		{
			Description: "Song exists",
			SongID:      mock.ValidSongID,
			Code:        http.StatusOK,
		},
		{
			Description: "Song exists, song lyrics not found",
			SongID:      mock.SongIDWithoutTextData,
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Song doesn't exists",
			SongID:      8923,
			Code:        http.StatusBadRequest,
		},
	}

	getSongDetailsHandler := handler.GetSongText(&mock.SongRepo{})

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {

			request := httptest.NewRequest("GET", "/api/v1/songs/{id}/lyrics", nil)

			request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprintf("%d", tc.SongID)})

			rr := httptest.NewRecorder()

			getSongDetailsHandler.ServeHTTP(rr, request)

			assert.Equal(t, tc.Code, rr.Code)
		})
	}
}
