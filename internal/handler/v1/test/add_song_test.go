package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amicie-monami/music-library/internal/domain/mock"
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/stretchr/testify/assert"
)

func TestAddSong(t *testing.T) {
	testCases := []struct {
		Description string
		ReqBody     any
		Code        int
	}{
		{
			Description: "Valid request body",
			ReqBody:     map[string]any{"Group": "Group", "Song": "Song"},
			Code:        http.StatusCreated,
		},
		{
			Description: "Group field is missing",
			ReqBody:     map[string]any{"Song": "Song"},
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Group field has zero value",
			ReqBody:     map[string]any{"Group": "", "Song": "Song"},
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Song field is missing",
			ReqBody:     map[string]any{"Group": "Group"},
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Song field has zero value",
			ReqBody:     map[string]any{"Group": "Group", "Song": ""},
			Code:        http.StatusBadRequest,
		},
		{
			Description: "Song and group fields have zero value",
			ReqBody:     map[string]any{"Group": "", "Song": ""},
			Code:        http.StatusBadRequest,
		},
	}

	addSongHandler := handler.AddSong(&mock.SongRepo{})

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			body, _ := json.Marshal(tc.ReqBody)

			rr := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/songs", bytes.NewBuffer(body))

			addSongHandler.ServeHTTP(rr, request)
			assert.Equal(t, tc.Code, rr.Code)
		})
	}
}
