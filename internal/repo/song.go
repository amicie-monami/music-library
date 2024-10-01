package repo

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/amicie-monami/music-library/pkg/myreflect"

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

func (r *Song) GetSongWithDetails(group string, title string) (*dto.SongWithDetails, error) {
	slog.Debug("get song", "group", group, "title", title)
	return r.getSongWithDetails(group, title)
}

func (r *Song) UpdateSong(song *model.Song) (int64, error) {
	slog.Debug("update song", "data", fmt.Sprintf("%+v", song))
	return r.updateSong(song)
}

func (r *Song) UpdateSongDetails(details *model.SongDetail) (int64, error) {
	slog.Debug("update song details", "data", fmt.Sprintf("%+v", details))
	return r.updateSongDetails(details)
}

func (r *Song) Delete(id int64) (int64, error) {
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
			"song_name",
		).
		Values(song.Group, song.Name).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id")

	sql, args := query.MustSql()
	return r.db.QueryRow(sql, args...).Scan(&song.ID)
}

func (r *Song) getSongs(aggregation map[string]any) ([]dto.SongWithDetails, error) {
	columns := getSongsBuildColumnNames(aggregation["fields"].(string))
	constraints, err := buildGetSongsWhereExpr(aggregation["filter"].(map[string]any))
	if err != nil {
		return nil, err
	}

	queryBuilder := squirrel.
		Select(columns...).
		From("songs").
		Join("song_details ON songs.id = song_details.song_id").
		Where(constraints).
		OrderBy("song_id").
		PlaceholderFormat(squirrel.Dollar)

	if aggregation["limit"] != "" {
		queryBuilder = queryBuilder.Limit(uint64(aggregation["limit"].(int64)))
	}

	if aggregation["offset"] != "" {
		queryBuilder = queryBuilder.Offset(uint64(aggregation["offset"].(int64)))
	}

	query, args := queryBuilder.MustSql()

	songs := make([]dto.SongWithDetails, 0)
	return songs, r.db.Select(&songs, query, args...)
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

func (r *Song) getSongWithDetails(group string, title string) (*dto.SongWithDetails, error) {
	var songWithDetails dto.SongWithDetails

	query, args := squirrel.
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
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	err := r.db.Get(&songWithDetails, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &songWithDetails, nil
}

func (r *Song) updateSong(song *model.Song) (int64, error) {
	if song == nil {
		return 0, nil
	}
	fields := myreflect.FillMapNotZeros(map[string]any{
		"group_name": song.Group,
		"song_name":  song.Name,
	})
	return updateRow(r.db, "songs", "id", song.ID, fields)
}

func (r *Song) updateSongDetails(details *model.SongDetail) (int64, error) {
	if details == nil {
		return 0, nil
	}
	fields := myreflect.FillMapNotZeros(map[string]any{
		"text":         details.Text,
		"link":         details.Link,
		"release_date": details.ReleaseDate,
	})
	return updateRow(r.db, "song_details", "song_id", details.SongID, fields)
}

func (r *Song) deleteById(id int64) (int64, error) {
	query, args := squirrel.
		Delete("songs").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func updateRow(db *sqlx.DB, table, pk_column string, pk any, fields map[string]any) (int64, error) {
	if len(fields) == 0 {
		return 0, nil
	}

	query := squirrel.Update(table).
		SetMap(fields).
		Where(squirrel.Eq{pk_column: pk}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args := query.MustSql()

	result, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
