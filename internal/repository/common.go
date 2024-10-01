package repository

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// updateRow Updates a row in the specified table.
// It takes
//   - db - database connection
//   - table - table name,
//   - pkEqConstraunt - a primary key equals constraint
//   - setMap - a map of column values (setMap) to be updated.
//
// The function filters out zero-value entries (nil, 0, empty string) from setMap to avoid updating columns with "zero" values.
// If no non-zero values are present in setMap, the function returns (0, nil) indicating no update was performed.
//
// Returns
//   - the number of affected rows (int64)
//   - an error if the query fails.
func updateRow(db *sqlx.DB, table string, pkEqCostraint squirrel.Eq, setMap map[string]any) (int64, error) {
	setMapWithoutZeros := make(map[string]any)

	for key, value := range setMap {
		if value != nil && value != 0 && value != "" {
			setMapWithoutZeros[key] = value
		}
	}

	if len(setMapWithoutZeros) == 0 {
		return 0, nil
	}

	query, args := squirrel.
		Update(table).
		SetMap(setMap).
		Where(pkEqCostraint).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// buildParamBasedAndConditions constructs a sql "and" condition [cond1 AND cond2 AND cond3...]
// based on the provided filter map and condition resolvers
//
// This function iterates over the filter map, where each key is a parameter name, and the corresponding value is a string
// representing the filter condition. For each key in the filter map, it finds a corresponding function in the
// conditionResolvers map, which is responsible for constructing the sql condition for that parameter
//
// Parameters:
//   - filter: A map where keys are parameter names (e.g., "song_name", "group") and values are their respective filter conditions
//   - conditionResolvers: A map where keys match the parameter names in the filter, and values are functions (conditionBuilderFunc)
//     that generate the SQL condition for the corresponding parameter
//
// Returns:
//   - squirrel.Sqlizer: A SQL condition builder object representing the combined "and" conditions
//   - error: An error if there was an issue constructing any condition
//   - nil if no conditions generated
func buildParamBasedAndConditions(params map[string]any, conditionResolvers map[string]conditionBuilderFunc) (squirrel.Sqlizer, error) {
	var conditions squirrel.And

	for param, builderFunc := range conditionResolvers {
		paramValue, ok := params[param].(string)
		if !ok || paramValue == "" {
			continue
		}

		condition, err := builderFunc(paramValue)
		if err != nil {
			return nil, err
		}

		if condition != nil {
			conditions = append(conditions, condition)
		}
	}

	if len(conditions) == 0 {
		return nil, nil
	}

	return conditions, nil
}

// buildILikeCondition constructs an SQL "ILIKE" condition [col ILIKE pattern1 OR col ILIKE pattern2...]
// for the given columnName using the provided pattern(s).
//
// The function takes a column name and a string of patterns. If multiple patterns are provided, separated by spaces, the function
// generates an "OR" condition, where each pattern is applied to the specified column using "ILIKE".
//
// Parameters:
//   - columnName: The name of the column to apply the "ILIKE" condition on.
//   - patterns: A string of patterns to match, where "*" is treated as a wildcard and spaces separate multiple patterns.
//
// Returns:
//   - squirrel.Sqlizer: A SQL "OR" condition with each pattern applied using the "ILIKE" operator.
//   - error: An error if the provided patterns string is empty.
func buildILikeCondition(columnName string, patterns string) (squirrel.Sqlizer, error) {
	if patterns == "" {
		return nil, fmt.Errorf("invalid value of %s param", columnName)
	}

	patterns = strings.ReplaceAll(patterns, "*", "%")
	parts := strings.Split(patterns, " ")

	orCondition := squirrel.Or{}
	for _, part := range parts {
		orCondition = append(orCondition, squirrel.ILike{columnName: part})
	}

	return orCondition, nil
}

// parseComparison converts a comparison string into the corresponding SQL operator.
//
//	Supported values are:
//	- "gt", "ge", "lt", "le", "ne", "eq"
func parseComparison(value string) string {
	switch value {
	case "gt":
		return ">"
	case "ge":
		return ">="
	case "lt":
		return "<"
	case "le":
		return "<="
	case "ne":
		return "!="
	case "eq":
		return "="
	default:
		return ""
	}
}
