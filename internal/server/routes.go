package server

import (
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/amicie-monami/music-library/internal/repo"
	"github.com/amicie-monami/music-library/internal/usecase"
	"github.com/amicie-monami/music-library/pkg/middleware"
	"github.com/gorilla/mux"
)

/*
	1.Получение данных библиотеки
		с фильтрацией по всем полям
		и пагинацией

	2.Получение текста песни
		с пагинацией по куплетам

	3.Удаление песни

	4.Изменение данных песни

	5. Добавление новой песни в формате

	6. /info
*/

func configureRouter(router *mux.Router, songRepo *repo.Song, songUsecase *usecase.Song) {
	//logs every request
	router.Use(middleware.Log)

	router.Handle("/api/v1/songs", handler.GetSongsData(songUsecase)).Methods("GET")

	//get the song text
	router.Handle("/api/v1/songs/{id}/lyrics", handler.GetSongText(songRepo)).Methods("GET")

	//delete the song
	router.Handle("/api/v1/songs/{id}", handler.DeleteSong(songRepo)).Methods("DELETE")

	//update data of the song
	router.Handle("/api/v1/songs/{id}", handler.UpdateSong(songRepo)).Methods("PATCH")

	//add a song to the library
	router.Handle("/api/v1/songs", handler.AddSong(songRepo)).Methods("POST")

	//get song details
	router.Handle("/api/v1/info", handler.GetSongDetails(songRepo)).Methods("GET")
}
