package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

// Song object adapter for database operations with songs tables
// (!!!) refactoring required. add context support
type Song struct {
	db *sqlx.DB
}

func NewSong(db *sqlx.DB) *Song {
	return &Song{db}
}

/// ------------ Interface ------------ ///

// Tx doesnt works, needs refactoring: e.g. Create a transactor object that can perform transactions
// (in this repository it is only needed when using update operators collectively)
func (r *Song) Tx(txActions func() error) error {
	tx := r.db.MustBegin()
	if err := txActions(); err != nil {
		tx.Rollback()
	}
	return tx.Commit()
}

func (r *Song) Create(song *model.Song) error {
	slog.Debug("create song", "data", fmt.Sprintf("%+v", song))

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

func (r *Song) GetSongs(aggregation map[string]any) ([]dto.SongWithDetails, error) {
	slog.Debug("get song", "aggregation data=", aggregation)

	columns := buildGetSongsColumnNames(aggregation["fields"].(string))
	whereExpr, err := buildGetSongsWhereExpr(aggregation["filter"].(map[string]any))
	if err != nil {
		return nil, err
	}

	queryBuilder := squirrel.
		Select(columns...).
		From("songs").
		Join("song_details ON songs.id = song_details.song_id").
		Where(whereExpr).
		OrderBy("song_id").
		PlaceholderFormat(squirrel.Dollar)

	if aggregation["limit"] != "" {
		queryBuilder = queryBuilder.Limit(uint64(aggregation["limit"].(int64)))
	} else {
		queryBuilder = queryBuilder.Limit(1000)
	}

	if aggregation["offset"] != "" {
		queryBuilder = queryBuilder.Offset(uint64(aggregation["offset"].(int64)))
	}

	query, args := queryBuilder.MustSql()

	songs := make([]dto.SongWithDetails, 0)
	return songs, r.db.Select(&songs, query, args...)
}

func (r *Song) GetSongText(id int64) (*string, error) {
	slog.Debug("get song text", "id", id)

	query, args := squirrel.
		Select("text").
		From("song_details").
		Where(squirrel.Eq{"song_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	songText := new(string)
	return songText, r.db.QueryRow(query, args...).Scan(songText)
}

func (r *Song) GetSongWithDetails(group string, title string) (*dto.SongWithDetails, error) {
	slog.Debug("get song", "group", group, "title", title)
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

func (r *Song) UpdateSong(song *model.Song) (int64, error) {
	slog.Debug("update song", "data", fmt.Sprintf("%+v", song))
	if song == nil {
		return 0, nil
	}

	table := "songs"
	primaryKeyEqauls := squirrel.Eq{"id": song.ID}

	setMap := map[string]any{
		"group_name": song.Group,
		"song_name":  song.Name,
	}

	return updateRow(r.db, table, primaryKeyEqauls, setMap)
}

func (r *Song) UpdateSongDetails(details *model.SongDetail) (int64, error) {
	slog.Debug("update song details", "data", fmt.Sprintf("%+v", details))

	table := "song_details"
	primaryKeyEqauls := squirrel.Eq{"song_id": details.SongID}

	setMap := map[string]any{
		"text":         details.Text,
		"link":         details.Link,
		"release_date": details.ReleaseDate,
	}

	return updateRow(r.db, table, primaryKeyEqauls, setMap)
}

func (r *Song) Delete(id int64) (int64, error) {
	slog.Debug("delete song", "id", id)

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

/// ------------ Helpers ------------ ///

func buildGetSongsColumnNames(columns string) []string {
	if columns == "" {
		return []string{
			"song_id",
			"group_name",
			"song_name",
			"release_date",
			"link",
			"text",
		}
	}
	return strings.Split(columns, " ")
}

// conditionBuilderFunc defines a function type that constructs SQL conditions.
type conditionBuilderFunc func(string) (squirrel.Sqlizer, error)

// buildGetSongsWhereExpr builds where expr
func buildGetSongsWhereExpr(filter map[string]any) (squirrel.Sqlizer, error) {
	conditionResolvers := map[string]conditionBuilderFunc{
		"song_id":      buildSongIDCondition,
		"song_name":    buildSongNameCondition,
		"groups":       buildGroupsCondition,
		"link":         buildLinkCondition,
		"release_date": buildReleaseDateCondition,
		"text":         buildSongTextCondition,
	}

	return buildParamBasedAndConditions(filter, conditionResolvers)
}

// buildSongIDCondition build comparison expression on song_id column [e.g. song_id >= 12]
func buildSongIDCondition(paramSongID string) (squirrel.Sqlizer, error) {
	if paramSongID == "" {
		return nil, fmt.Errorf("invalid value of song_id param")
	}

	songIDConstraint := strings.Split(paramSongID, " ")

	//assume that parameter hasn't comparison operator
	comparison := "="
	songID := songIDConstraint[0]

	//if split returns two elements, that parameter must contain a comparision operator.
	//Change default "=" on passed operator
	if len(songIDConstraint) == 2 {
		comparison = parseComparison(songIDConstraint[0])
		if comparison == "" {
			return nil, fmt.Errorf("invalid comparison operator in sond_id param")
		}
		songID = songIDConstraint[1]

	} else if len(songIDConstraint) > 2 {
		and := squirrel.And{}
		for idx := 0; idx < len(songIDConstraint); idx += 2 {
			comparison = parseComparison(songIDConstraint[idx])
			if comparison == "" {
				return nil, fmt.Errorf("invalid comparison operator in sond_id param: %s, %v", songIDConstraint[idx], songIDConstraint)
			}
			songID = songIDConstraint[idx+1]
			and = append(and, squirrel.Expr(fmt.Sprintf("song_id %s ?", comparison), songID))
		}

		return and, nil
	}

	//validation
	if _, err := strconv.ParseInt(songID, 10, 64); err != nil {
		return nil, fmt.Errorf("invalid value in song_id param")
	}

	//build comparison expression fot the song_id column
	return squirrel.Expr(fmt.Sprintf("song_id %s ?", comparison), songID), nil
}

// buildSongNameCondition builds an ILIKE expression (or multiple ILIKEs
// associated with the OR operator) for the "song_name" column.
func buildSongNameCondition(paramSongName string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("song_name", paramSongName)
}

// buildGroupsCondition builds equals expression (or muliplie equals associated with the OR operator)
// for the "group_name" column
func buildGroupsCondition(paramGroups string) (squirrel.Sqlizer, error) {
	if paramGroups == "" {
		return nil, fmt.Errorf("invalid value of groups param")
	}
	groups := strings.Split(paramGroups, " ")
	for idx := range groups {
		groups[idx] = strings.ReplaceAll(groups[idx], "_", " ")
	}
	return squirrel.Eq{"group_name": groups}, nil
}

// buildSongNameCondition builds an ILIKE expression (or multiple ILIKEs
// associated with the OR operator) for the "link" column.
func buildLinkCondition(paramLink string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("link", paramLink)
}

// buildReleaseDateCondition builds constraints for the "release_date" column depending on the content of the paramter
//   - Comparision expr (if the param is [comparison operator]+[date] kind)
//   - Equals expr (if the parameter has only one date)
//   - Between expr(if the parameteris of the [date]-[date] kind)
func buildReleaseDateCondition(paramReleaseDate string) (squirrel.Sqlizer, error) {
	if paramReleaseDate == "" {
		return nil, fmt.Errorf("invalid value of release_date param")
	}

	layout := "01.01.2006"
	dateConstraint := strings.Split(paramReleaseDate, " ")

	//release date parameter has param=[operator]+[date] kind
	if len(dateConstraint) == 2 {
		date, err := time.Parse(layout, dateConstraint[1])
		if err != nil {
			return nil, err
		}

		comparison := parseComparison(dateConstraint[0])
		if comparison == "" {
			return nil, fmt.Errorf("invalid comparison operator %s", dateConstraint[0])
		}

		return squirrel.Expr(fmt.Sprintf("release_date %s ?", comparison), date), nil
	}

	dates := strings.Split(paramReleaseDate, "-")
	//release date parameter has param=[date] kind
	if len(dates) == 1 {
		date, err := time.Parse(layout, dates[0])
		if err != nil {
			return nil, err
		}
		return squirrel.Eq{"release_date": date}, nil
	}

	//release date parameter has param=[date]-[date] kind
	if len(dates) == 2 {
		startDate, err := time.Parse(layout, dates[0])
		if err != nil {
			return nil, err
		}

		endDate, err := time.Parse(layout, dates[1])
		if err != nil {
			return nil, err
		}
		return squirrel.Expr("release_date BETWEEN ? AND ?", startDate, endDate), nil
	}

	return nil, fmt.Errorf("invalid value of release_date param=%s", paramReleaseDate)
}

// buildSongTextCondition builds an ILIKE expression (or multiple ILIKEs
// associated with the OR operator) for the "text" column.
func buildSongTextCondition(paramSongText string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("text", paramSongText)
}
