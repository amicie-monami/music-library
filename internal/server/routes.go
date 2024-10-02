package server

import (
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/amicie-monami/music-library/internal/repository"
	"github.com/amicie-monami/music-library/pkg/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func configureRouter(router *mux.Router, songRepo *repository.Song) {

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.Handle("/api/v1/songs", middleware.Log(handler.GetSongsData(songRepo))).Methods("GET")

	router.Handle("/api/v1/songs/{id}/lyrics", middleware.Log(handler.GetSongText(songRepo))).Methods("GET")

	router.Handle("/api/v1/songs/{id}", middleware.Log(handler.DeleteSong(songRepo))).Methods("DELETE")

	router.Handle("/api/v1/songs/{id}", middleware.Log(handler.UpdateSong(songRepo))).Methods("PATCH")

	router.Handle("/api/v1/songs", middleware.Log(handler.AddSong(songRepo))).Methods("POST")

	router.Handle("/api/v1/info", middleware.Log(handler.GetSongDetails(songRepo))).Methods("GET")
}
