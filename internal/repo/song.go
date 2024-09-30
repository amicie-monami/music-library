package repo

import (
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
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
	slog.Debug("create song", "data", fmt.Sprintf("%+v", song))
	return r.create(song)
}

func (r *Song) GetSongs(aggregation map[string]any) ([]dto.SongWithDetails, error) {
	slog.Debug("get song", "aggregation data=", aggregation)
	return r.getSongs(aggregation)
}

func (r *Song) GetSongText(id int64) (*string, error) {
	slog.Debug("get song text", "id", id)
	return r.getSongTextById(id)
}

func (r *Song) GetSongDetails(group string, title string) (*model.SongDetail, error) {
	slog.Debug("get song", "group", group, "title", title)
	return r.getSongDetails(group, title)
}

func (r *Song) UpdateSong(song *model.Song) error {
	slog.Debug("update song", "data", fmt.Sprintf("%+v", song))
	return r.updateSong(song)
}

func (r *Song) UpdateSongDetails(details *model.SongDetail) error {
	slog.Debug("update song details", "data", fmt.Sprintf("%+v", details))
	return r.updateSongDetails(details)
}

func (r *Song) Delete(id int64) error {
	slog.Debug("delete song", "id", id)
	return r.deleteById(id)
}

func (r *Song) Tx(txActions func() error) error {
	return txActions()
}

func (r *Song) create(song *model.Song) error {
	query := squirrel.
		Insert("songs").
		Columns(
			"group_name",
			"song_title",
		).
		Values(song.Group, song.Name).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id")

	sql, args := query.MustSql()
	return r.db.QueryRow(sql, args...).Scan(&song.ID)
}

func (r *Song) getSongs(aggregation map[string]any) ([]dto.SongWithDetails, error) {
	columns := getSongsBuildColumnNames(aggregation["fields"].(string))
	confitions, err := buildGetSongsWhereExpr(aggregation["filter"].(map[string]any))
	if err != nil {
		return nil, err
	}

	query := squirrel.
		Select(columns...).
		From("songs").
		Join("song_details ON songs.id = song_details.song_id").
		Where(confitions).
		OrderBy("song_id")
		// Limit(uint64(aggregation["limit"].(int64))).
		// Offset(uint64(aggregation["offset"].(int64)))

	sql, args := query.PlaceholderFormat(squirrel.Dollar).MustSql()
	fmt.Println(sql)

	songs := make([]dto.SongWithDetails, 0)
	return songs, r.db.Select(&songs, sql, args...)
}

func (r *Song) getSongTextById(id int64) (*string, error) {
	query := squirrel.
		Select("text").
		From("song_details").
		Where(squirrel.Eq{"song_id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()
	songText := new(string)
	return songText, r.db.QueryRow(sql, args...).Scan(songText)
}

func (r *Song) getSongDetails(group string, title string) (*model.SongDetail, error) {
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
			"songs.song_name":  title,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()
	details := new(model.SongDetail)
	return details, r.db.Get(details, sql, args...)
}

func (r *Song) updateSong(song *model.Song) error {
	if song == nil {
		return nil
	}
	fields := reflect.FillMapNotZeros(map[string]any{
		"group_name": song.Group,
		"song_name":  song.Name,
	})
	return updateRow(r.db, "songs", "id", song.ID, fields)
}

func (r *Song) updateSongDetails(details *model.SongDetail) error {
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

func (r *Song) deleteById(id int64) error {
	query := squirrel.
		Delete("songs").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()
	_, err := r.db.Exec(sql, args...)
	return err
}

func updateRow(db *sqlx.DB, table, pk_column string, pk any, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	query := squirrel.Update(table).
		SetMap(fields).
		Where(squirrel.Eq{pk_column: pk}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()

	result, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}
