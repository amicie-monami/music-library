package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

// dbContext describes the database context.
// Can take the values sqlx.DB or sqlx.TX
type dbContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Song object adapter for database operations with songs tables
type Song struct {
	db dbContext
}

func NewSong(db dbContext) *Song {
	return &Song{db}
}

/// ------------ Interface ------------ ///

// Tx execute the transaction, the action of which described in the txActions
func (r *Song) Tx(ctx context.Context, txActions func() error) error {
	db, ok := r.db.(*sqlx.DB)
	if !ok {
		debugMessage := fmt.Sprintf("execute the tx, unsupport database type: %v", reflect.TypeOf(r.db))
		return dto.NewError(500, "internal server error", "song.Tx", nil, debugMessage)
	}

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		debugMessage := fmt.Errorf("failed to begin the transaction: %s", err)
		return dto.NewError(500, "internal server error", "song.Tx", nil, debugMessage)
	}

	//change the context of the repository database to a transaction
	r.db = tx
	defer func() {
		//after transactions are executed, return the database context to the instance
		r.db = db
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	//make a payload
	if err = txActions(); err != nil {
		return err
	}

	//commits the changes
	if err = tx.Commit(); err != nil {
		debugMessage := fmt.Errorf("failed to commit the transaction: %s", err)
		return dto.NewError(500, "internal server error", "song.Tx", nil, debugMessage)
	}

	return nil
}

func (r *Song) Create(ctx context.Context, song *model.Song) error {
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
	if err := r.db.QueryRowContext(ctx, sql, args...).Scan(&song.ID); err != nil {
		return wrapQueryExecError("song.Create", err)
	}

	return nil
}

func (r *Song) GetSongs(ctx context.Context, aggregation map[string]any) ([]*dto.SongWithDetails, error) {
	slog.Debug("get song", "aggregation data=", aggregation)

	//setup aggragation filters
	columns := buildGetSongsColumnNames(aggregation["fields"].(string))
	whereExpr, err := buildGetSongsWhereExpr(aggregation["filter"].(map[string]any))
	if err != nil {
		return nil, err
	}

	//mark the body of the sql query
	queryBuilder := squirrel.
		Select(columns...).
		From("songs").
		Join("song_details ON songs.id = song_details.song_id").
		Where(whereExpr).
		OrderBy("song_id").
		PlaceholderFormat(squirrel.Dollar)

	//setup pagination filters
	if aggregation["limit"] != "" {
		queryBuilder = queryBuilder.Limit(uint64(aggregation["limit"].(int64)))
	} else {
		queryBuilder = queryBuilder.Limit(1000)
	}
	if aggregation["offset"] != "" {
		queryBuilder = queryBuilder.Offset(uint64(aggregation["offset"].(int64)))
	}

	//build sql query
	query, args := queryBuilder.MustSql()

	//execution
	songs := make([]*dto.SongWithDetails, 0)
	if err := r.db.SelectContext(ctx, &songs, query, args...); err != nil {
		return nil, wrapQueryExecError("song.GetSongs", err)
	}

	fmt.Println(query)

	return songs, nil
}

func (r *Song) GetSongText(ctx context.Context, id int64) (*string, error) {
	slog.Debug("get song text", "id", id)

	query, args := squirrel.
		Select("text").
		From("song_details").
		Where(squirrel.Eq{"song_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	songText := new(*string)
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(songText); err != nil {

		if err == sql.ErrNoRows {
			message := "could not found the song"
			details := fmt.Sprintf("id=%d", id)
			return nil, dto.NewError(400, message, "song.GetSongText", details, nil)
		}

		return nil, wrapQueryExecError("song.GetSongs", err)
	}

	return *songText, nil
}

func (r *Song) GetSongWithDetails(ctx context.Context, group string, title string) (*dto.SongWithDetails, error) {
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

	err := r.db.GetContext(ctx, &songWithDetails, query, args...)
	if err != nil {

		if err == sql.ErrNoRows {
			message := "could not found the song"
			details := fmt.Sprintf("song=%s, group=%s", title, group)
			return nil, dto.NewError(400, message, "song.GetSongText", details, nil)
		}

		return nil, wrapQueryExecError("song.GetSongWithDetails", err)
	}

	return &songWithDetails, nil
}

func (r *Song) UpdateSong(ctx context.Context, song *model.Song) error {
	slog.Debug("update song", "data", fmt.Sprintf("%+v", song))
	if song == nil {
		return dto.NewError(500, "internal server error", "song.UpdateSong", nil, "song is nil")
	}

	table := "songs"
	primaryKeyEqauls := squirrel.Eq{"id": song.ID}

	setMap := map[string]any{
		"group_name": song.Group,
		"song_name":  song.Name,
	}

	affectedCount, err := updateRowContext(ctx, r.db, table, primaryKeyEqauls, setMap)
	if err != nil {
		return wrapQueryExecError("song.UpdateSong", err)
	}

	if affectedCount == 0 {
		details := fmt.Sprintf("id=%d", song.ID)
		return dto.NewError(400, "song not found", "song.UpdateSong", details, nil)
	}

	return nil
}

func (r *Song) UpdateSongDetails(ctx context.Context, details *model.SongDetail) error {
	slog.Debug("update song details", "data", fmt.Sprintf("%+v", details))

	table := "song_details"
	primaryKeyEqauls := squirrel.Eq{"song_id": details.SongID}

	setMap := map[string]any{
		"text":         details.Text,
		"link":         details.Link,
		"release_date": details.ReleaseDate,
	}

	affectedCount, err := updateRowContext(ctx, r.db, table, primaryKeyEqauls, setMap)
	if err != nil {
		return wrapQueryExecError("song.UpdateSong", err)
	}

	if affectedCount == 0 {
		details := fmt.Sprintf("id=%d", details.SongID)
		return dto.NewError(400, "song not found", "song.UpdateSong", details, nil)
	}

	return nil
}

func (r *Song) Delete(ctx context.Context, id int64) error {
	slog.Debug("delete song", "id", id)

	query, args := squirrel.
		Delete("songs").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return wrapQueryExecError("song.Delete", err)
	}

	affectedCount, err := result.RowsAffected()
	if err != nil {
		return wrapQueryExecError("song.Delete", err)
	}

	if affectedCount == 0 {
		details := fmt.Sprintf("id=%d", id)
		return dto.NewError(400, "song not found", "song.Delete", details, nil)
	}

	return nil
}

/// ------------ Helpers ------------ ///

func wrapQueryExecError(source string, err error) *dto.Error {
	debugMsg := fmt.Errorf("database error: %v", err)
	return dto.NewError(500, "internal server error", source, nil, debugMsg)
}

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
		return nil, dto.NewError(400, "missing value in song_id param ", "buildSongIDCondition", nil, nil)
	}

	var err error
	songIDConstraint := strings.Split(paramSongID, " ")

	//assume that parameter hasn't comparison operator
	comparison := "="
	songID := songIDConstraint[0]

	//if split returns two elements, that parameter must contain a comparision operator.
	//Change default "=" on param's operator
	if len(songIDConstraint) == 2 {
		comparison, err = parseComparisonOperatorInSongIDParam(songIDConstraint[0])
		if err != nil {
			return nil, err
		}

		songID = songIDConstraint[1]

		//if param has 2 and more values
	} else if len(songIDConstraint) > 2 {
		and := squirrel.And{}
		//range for pairs [operator constraint, operator constraint]
		for idx := 0; idx < len(songIDConstraint); idx += 2 {
			comparison, err = parseComparisonOperatorInSongIDParam(songIDConstraint[idx])
			if err != nil {
				return nil, err
			}

			songID = songIDConstraint[idx+1]
			and = append(and, squirrel.Expr(fmt.Sprintf("song_id %s ?", comparison), songID))
		}

		return and, nil
	}

	//validation
	if _, err := strconv.ParseInt(songID, 10, 64); err != nil {
		return nil, dto.NewError(400, "song_id comstraint must be a num", "buildSongIDCondition", songID, nil)
	}

	//build comparison expression fot the song_id column
	return squirrel.Expr(fmt.Sprintf("song_id %s ?", comparison), songID), nil
}

func parseComparisonOperatorInSongIDParam(operator string) (string, error) {
	comparison := parseComparison(operator)
	if comparison == "" {
		return "", dto.NewError(400, "invalid comparison operator in song_id param", "buildSongIDCondition", operator, nil)
	}
	return comparison, nil
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
		return nil, dto.NewError(400, "missing value in groups param ", "buildSongIDCondition", nil, nil)
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
		return nil, dto.NewError(400, "missing value in release_date param ", "buildSongIDCondition", nil, nil)
	}

	dateConstraint := strings.Split(paramReleaseDate, " ")

	//release date parameter has param=[operator]+[date] kind
	if len(dateConstraint) == 2 {
		date, err := parseDateInReleaseDateParam(dateConstraint[1])
		if err != nil {
			return nil, err
		}

		comparison, err := parseComparisonOperatorInReleaseDateParam(dateConstraint[0])
		if err != nil {
			return nil, err
		}

		return squirrel.Expr(fmt.Sprintf("release_date %s ?", comparison), date), nil
	}

	dates := strings.Split(paramReleaseDate, "-")
	//release date parameter has param=[date] kind
	if len(dates) == 1 {

		date, err := parseDateInReleaseDateParam(dates[0])
		if err != nil {
			return nil, err
		}

		return squirrel.Eq{"release_date": date}, nil
	}

	//release date parameter has param=[date]-[date] kind
	if len(dates) == 2 {

		startDate, err := parseDateInReleaseDateParam(dates[0])
		if err != nil {
			return nil, err
		}

		endDate, err := parseDateInReleaseDateParam(dates[1])
		if err != nil {
			return nil, err
		}

		return squirrel.Expr("release_date BETWEEN ? AND ?", startDate, endDate), nil
	}

	return nil, dto.NewError(400, "failed to parse release_date param", "buildReleaseDateCondition", paramReleaseDate, nil)
}

func parseDateInReleaseDateParam(dateStr string) (time.Time, error) {
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return time.Time{}, dto.NewError(400, "could not parse the release_date, expected format was dd.mm.yyyy", "buildReleaseDateCondition", dateStr, nil)
	}
	return date, nil
}

func parseComparisonOperatorInReleaseDateParam(operator string) (string, error) {
	comparison := parseComparison(operator)
	if comparison == "" {
		return "", dto.NewError(400, "invalid comparison operator in release_date param", "buildReleaseDateCondition", operator, nil)
	}
	return comparison, nil
}

// buildSongTextCondition builds an ILIKE expression (or multiple ILIKEs
// associated with the OR operator) for the "text" column.
func buildSongTextCondition(paramSongText string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("text", paramSongText)
}
