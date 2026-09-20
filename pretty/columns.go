package pretty

import (
	"reflect"
	"sort"
)

func collectColumns(rows []map[string]any) []string {
	seen := map[string]struct{}{}
	for _, row := range rows {
		for key := range row {
			seen[key] = struct{}{}
		}
	}
	columns := make([]string, 0, len(seen))
	for key := range seen {
		columns = append(columns, key)
	}
	return columns
}

func extractConstants(rows []map[string]any, columns []string) (map[string]any, []string) {
	constants := map[string]any{}
	if len(rows) < 2 {
		return constants, columns
	}
	variable := columns[:0]
	for _, column := range columns {
		first, exists := rows[0][column]
		if !exists || first == nil || !isScalar(first) {
			variable = append(variable, column)
			continue
		}
		constant := true
		for _, row := range rows[1:] {
			value, exists := row[column]
			if !exists || !reflect.DeepEqual(value, first) {
				constant = false
				break
			}
		}
		if constant {
			constants[column] = first
		} else {
			variable = append(variable, column)
		}
	}
	return constants, variable
}

func isScalar(value any) bool {
	switch value.(type) {
	case []any, map[string]any:
		return false
	default:
		return true
	}
}

func sortColumns(rows []map[string]any, columns []string) {
	sort.Slice(columns, func(i, j int) bool {
		left, right := columns[i], columns[j]
		leftNested := columnIsNested(rows, left)
		rightNested := columnIsNested(rows, right)
		if leftNested != rightNested {
			return !leftNested
		}
		return left < right
	})
}

func columnIsNested(rows []map[string]any, column string) bool {
	for _, row := range rows {
		switch row[column].(type) {
		case []any, map[string]any:
			return true
		}
	}
	return false
}
