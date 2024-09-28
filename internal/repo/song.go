package repo

import (
	"github.com/Masterminds/squirrel"
	"github.com/amicie-monami/music-library/internal/model"
	"github.com/jmoiron/sqlx"
)

type Song struct {
	db *sqlx.DB
}

func NewSong(db *sqlx.DB) *Song {
	return &Song{db}
}

func (r *Song) Create(song *model.Song) error {
	sql, args := squirrel.Insert("songs").Columns("group_name", "title").Values(song.Group, song.Title).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").MustSql()
	return r.db.QueryRow(sql, args...).Scan(&song.ID)
}

func (r *Song) Delete(id int64) error {
	sql, args := squirrel.Delete("songs").Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).MustSql()
	_, err := r.db.Exec(sql, args...)
	return err
}
