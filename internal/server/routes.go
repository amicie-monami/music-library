package server

import (
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/amicie-monami/music-library/internal/repo"
	"github.com/amicie-monami/music-library/pkg/middleware"
	"github.com/gorilla/mux"
)

/*
- Получение данных библиотеки
	с фильтрацией по всем полям
	и пагинацией
- Получение текста песни
	с пагинацией по куплетам
- Удаление песни
- Изменение данных песни
*/

func configureRouter(router *mux.Router, songRepo *repo.Song) {
	//logs every request
	router.Use(middleware.Log)

	//get library data
	router.Handle("/api/v1/songs", handler.GetData()).Methods("GET")

	//get the song text
	router.Handle("/api/v1/songs/{id}", handler.GetSongText()).Methods("GET")

	//delete the song
	router.Handle("/api/v1/songs/{id}", handler.DeleteSong()).Methods("DELETE")

	//update data of the song
	router.Handle("/api/v1/songs/{id}", handler.UpdateSong()).Methods("PATCH")

	//add a song to the library
	router.Handle("/api/v1/songs", handler.AddSong(songRepo)).Methods("POST")
}
