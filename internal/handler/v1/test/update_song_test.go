package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amicie-monami/music-library/internal/domain/mock"
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSong(t *testing.T) {
	testCases := []struct {
		Description string
		ReqBody     any
		SongID      int64
		Code        int
	}{
		{
			Description: "Valid request body",
			ReqBody:     map[string]any{"Group": "Group", "release_date": "01.02.2022"},
			SongID:      mock.ValidSongID,
			Code:        http.StatusOK,
		},
		{
			Description: "Empty request body",
			ReqBody:     map[string]any{},
			SongID:      mock.ValidSongID,
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Invalid release date fields",
			ReqBody:     map[string]any{"release_date": "2022.03.05"},
			SongID:      mock.ValidSongID,
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Song not found",
			ReqBody:     map[string]any{"release_date": "2022.03.05"},
			SongID:      9090,
			Code:        http.StatusBadRequest,
		},
	}

	addSongHandler := handler.UpdateSong(&mock.SongRepo{})

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			body, _ := json.Marshal(tc.ReqBody)

			request := httptest.NewRequest("GET", "/api/songs/{id}", bytes.NewBuffer(body))

			request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprintf("%d", tc.SongID)})

			rr := httptest.NewRecorder()

			addSongHandler.ServeHTTP(rr, request)

			assert.Equal(t, tc.Code, rr.Code)
		})
	}
}
