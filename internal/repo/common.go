package repo

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

// conditionBuilderFunc defines a function type that constructs SQL conditions.
type conditionBuilderFunc func(string) (squirrel.Sqlizer, error)

// buildParamBasedAndConditions constructs a sql "and" condition based on the provided filter map and condition resolvers
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
func buildParamBasedAndConditions(filter map[string]any, conditionResolvers map[string]conditionBuilderFunc) (squirrel.Sqlizer, error) {
	var conditions squirrel.And

	for param, builderFunc := range conditionResolvers {
		paramValue, ok := filter[param].(string)
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

// buildILikeCondition constructs an SQL "ILIKE" condition for the given columnName using the provided pattern(s).
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
		return nil, fmt.Errorf("invalid value of groups param")
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
