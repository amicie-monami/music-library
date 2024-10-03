package handler

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

func TestDeleteSong(t *testing.T) {
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
			Description: "Song doesn't exist",
			SongID:      489,
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid song id",
			SongID:      -1,
			Code:        http.StatusBadRequest,
		},
	}

	deleteSongHandler := handler.DeleteSong(&mock.SongRepo{})

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {

			rr := httptest.NewRecorder()

			request := httptest.NewRequest("DELETE", "/api/v1/songs/id", nil)

			request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprintf("%d", tc.SongID)})

			deleteSongHandler.ServeHTTP(rr, request)

			assert.Equal(t, rr.Code, tc.Code)
		})

	}
}
