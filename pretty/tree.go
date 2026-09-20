package pretty

import (
	"fmt"
	"io"
	"strings"
)

func requiresTree(value any) bool {
	switch data := value.(type) {
	case map[string]any:
		for _, candidate := range data {
			switch nested := candidate.(type) {
			case map[string]any:
				return true
			case []any:
				if rows, ok := recordCollection(nested); ok {
					if !rowsAreFlat(rows) {
						return true
					}
					continue
				}
				for _, item := range nested {
					switch item.(type) {
					case map[string]any, []any:
						return true
					}
				}
			}
		}
	case []any:
		if rows, ok := recordCollection(data); ok {
			return !rowsAreFlat(rows)
		}
	}
	return false
}

func rowsAreFlat(rows []map[string]any) bool {
	for _, row := range rows {
		for _, value := range row {
			switch nested := value.(type) {
			case map[string]any:
				return false
			case []any:
				for _, item := range nested {
					if !isScalar(item) {
						return false
					}
				}
			}
		}
	}
	return true
}

func renderTree(w io.Writer, value any) error {
	var out strings.Builder
	switch data := value.(type) {
	case map[string]any:
		if err := renderTreeObject(&out, data, ""); err != nil {
			return err
		}
	case []any:
		rows, _ := recordCollection(data)
		out.WriteString("Items\n")
		for i, row := range rows {
			if err := renderTreeNode(&out, fmt.Sprintf("Item %d", i+1), row, "", i == len(rows)-1); err != nil {
				return err
			}
		}
	}
	_, err := io.WriteString(w, out.String())
	return err
}
