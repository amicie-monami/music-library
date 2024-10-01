package server

import (
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/amicie-monami/music-library/internal/repository"
	"github.com/amicie-monami/music-library/pkg/middleware"
	"github.com/gorilla/mux"
)

func configureRouter(router *mux.Router, songRepo *repository.Song) {
	//logs every request
	router.Use(middleware.Log)

	router.Handle("/api/v1/songs", handler.GetSongsData(songRepo)).Methods("GET")

	router.Handle("/api/v1/songs/{id}/lyrics", handler.GetSongText(songRepo)).Methods("GET")

	router.Handle("/api/v1/songs/{id}", handler.DeleteSong(songRepo)).Methods("DELETE")

	router.Handle("/api/v1/songs/{id}", handler.UpdateSong(songRepo)).Methods("PATCH")

	router.Handle("/api/v1/songs", handler.AddSong(songRepo)).Methods("POST")

	router.Handle("/api/v1/info", handler.GetSongDetails(songRepo)).Methods("GET")
}
