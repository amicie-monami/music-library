package usecase

import (
	"github.com/amicie-monami/music-library/internal/model"
	"github.com/amicie-monami/music-library/internal/repo"
)

type Song struct {
	repo *repo.Song
}

func New(repo *repo.Song) *Song {
	return &Song{repo}
}

func (u *Song) GetSongsData(aggregation map[string]any) ([]model.SongWithDetailsDTO, error) {
	return u.repo.GetSongs(aggregation)
}
