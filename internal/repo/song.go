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

func (r *Song) Create(song *model.Song) error {
	sql, args := squirrel.Insert("songs").Columns("group_name", "song_title").Values(song.Group, song.Title).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").MustSql()
	return r.db.QueryRow(sql, args...).Scan(&song.ID)
}

func (r *Song) Delete(id int64) error {
	sql, args := squirrel.Delete("songs").Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).MustSql()
	_, err := r.db.Exec(sql, args...)
	return err
}

func (r *Song) Update(song *model.Song, details *model.SongDetail) error {
	slog.Debug("update", "data", fmt.Sprintf("song=%+v, details=%+v", song, details))

	if song != nil {
		fields := reflect.FillMapNotZeros(map[string]any{
			"group_name": song.Group,
			"song_title": song.Title,
		})
		if err := updateRow(r.db, "songs", "id", song.ID, fields); err != nil {
			return err
		}
	}

	if details != nil {
		fields := reflect.FillMapNotZeros(map[string]any{
			"text":         details.Text,
			"link":         details.Link,
			"release_date": details.ReleaseDate,
		})
		if err := updateRow(r.db, "song_details", "song_id", details.SongID, fields); err != nil {
			return err
		}
	}

	return nil
}
