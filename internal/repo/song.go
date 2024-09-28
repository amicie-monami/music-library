package repo

import (
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/amicie-monami/music-library/internal/model"
	"github.com/amicie-monami/music-library/pkg/reflect"
	"github.com/jmoiron/sqlx"
)

type Song struct {
	db *sqlx.DB
}

func NewSong(db *sqlx.DB) *Song {
	return &Song{db}
}

func (r *Song) Tx(txActions func() error) error {
	return txActions()
}

func (r *Song) Create(song *model.Song) error {
	slog.Debug("create song", "data", fmt.Sprintf("%+v", song))
	query := squirrel.
		Insert("songs").
		Columns(
			"group_name",
			"song_title",
		).
		Values(song.Group, song.Title).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id")

	sql, args := query.MustSql()
	return r.db.QueryRow(sql, args...).Scan(&song.ID)
}

func (r *Song) GetSongText(id int64) (*string, error) {
	slog.Debug("get song text", "id", id)
	query := squirrel.
		Select("text").
		From("song_details").
		Where(squirrel.Eq{"song_id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()
	songText := new(string)
	return songText, r.db.QueryRow(sql, args...).Scan(songText)
}

func (r *Song) GetSongDetails(group string, title string) (*model.SongDetail, error) {
	slog.Debug("get song", "group", group, "title", title)
	query := squirrel.
		Select(
			"song_id",
			"release_date",
			"link",
			"text",
		).
		From("song_details").
		Join("songs on songs.id = song_id").
		Where(squirrel.Eq{
			"songs.group_name": group,
			"songs.song_title": title,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()
	details := new(model.SongDetail)
	return details, r.db.Get(details, sql, args...)
}

func (r *Song) UpdateSong(song *model.Song) error {
	slog.Debug("update song", "data", fmt.Sprintf("%+v", song))
	if song == nil {
		return nil
	}
	fields := reflect.FillMapNotZeros(map[string]any{
		"group_name": song.Group,
		"song_title": song.Title,
	})
	return updateRow(r.db, "songs", "id", song.ID, fields)
}

func (r *Song) UpdateSongDetails(details *model.SongDetail) error {
	slog.Debug("update song details", "data", fmt.Sprintf("%+v", details))
	if details == nil {
		return nil
	}
	fields := reflect.FillMapNotZeros(map[string]any{
		"text":         details.Text,
		"link":         details.Link,
		"release_date": details.ReleaseDate,
	})
	return updateRow(r.db, "song_details", "song_id", details.SongID, fields)
}

func (r *Song) Delete(id int64) error {
	slog.Debug("delete song", "id", id)
	query := squirrel.
		Delete("songs").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()
	_, err := r.db.Exec(sql, args...)
	return err
}
