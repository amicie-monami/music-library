package repo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
)

func getSongsBuildColumnNames(columns string) []string {
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

// buildGetSongsWhereExpr ...
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

func buildSongIDCondition(paramSongID string) (squirrel.Sqlizer, error) {
	if paramSongID == "" {
		return nil, fmt.Errorf("invalid value of song_id param")
	}

	songIDConstraint := strings.Split(paramSongID, " ")

	comparison := "="
	songID := songIDConstraint[0]

	if len(songIDConstraint) == 2 {
		comparison = parseComparison(songIDConstraint[0])
		if comparison == "" {
			return nil, fmt.Errorf("invalid comparison operator in sond_id param")
		}
		songID = songIDConstraint[1]
	}

	if _, err := strconv.ParseInt(songID, 10, 64); err != nil {
		return nil, fmt.Errorf("invalid value in song_id param")
	}

	return squirrel.Expr(fmt.Sprintf("song_id %s ?", comparison), songID), nil
}

func buildSongNameCondition(paramSongName string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("song_name", paramSongName)
}

func buildSongTextCondition(paramSongText string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("text", paramSongText)
}

func buildGroupsCondition(paramGroups string) (squirrel.Sqlizer, error) {
	if paramGroups == "" {
		return nil, fmt.Errorf("invalid value of groups param")
	}
	groups := strings.Split(paramGroups, " ")
	return squirrel.Eq{"group_name": groups}, nil
}

func buildLinkCondition(paramLink string) (squirrel.Sqlizer, error) {
	return buildILikeCondition("link", paramLink)
}

func buildReleaseDateCondition(paramReleaseDate string) (squirrel.Sqlizer, error) {
	if paramReleaseDate == "" {
		return nil, fmt.Errorf("invalid value of release_date param")
	}

	layout := "01.01.2006"
	dateConstraint := strings.Split(paramReleaseDate, " ")

	// if param has comparison operator
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
	if len(dates) == 1 {
		date, err := time.Parse(layout, dates[0])
		if err != nil {
			return nil, err
		}
		return squirrel.Eq{"release_date": date}, nil
	}

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

	return nil, nil
}
