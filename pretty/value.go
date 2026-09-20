package pretty

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func writeMetadata(w io.Writer, metadata map[string]any, byteColumns map[string]bool) error {
	_, err := fmt.Fprintln(w, formatMetadata(metadata, byteColumns))
	return err
}

func formatMetadata(metadata map[string]any, byteColumns map[string]bool) string {
	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, key := range keys {
		parts[i] = header(key) + "=" + cell(key, metadata[key], byteColumns)
	}
	return strings.Join(parts, "  ")
}

func cell(column string, value any, byteColumns map[string]bool) string {
	if byteColumns[column] {
		if n, ok := number(value); ok {
			return humanBytes(n)
		}
	}
	var result string
	switch data := value.(type) {
	case nil:
		return ""
	case string:
		result = data
	case bool:
		result = strconv.FormatBool(data)
	case []any:
		parts := make([]string, len(data))
		for i, item := range data {
			if !isScalar(item) {
				result = compactJSON(value)
				break
			}
			parts[i] = fmt.Sprint(item)
		}
		if result == "" {
			result = strings.Join(parts, ", ")
		}
	case map[string]any:
		result = compactJSON(value)
	default:
		result = fmt.Sprint(value)
	}
	result = strings.NewReplacer("\n", " ", "\r", " ", "\t", " ").Replace(result)
	return result
}

func tableCell(column string, value any, byteColumns map[string]bool) any {
	if !byteColumns[column] {
		if _, ok := number(value); ok {
			return value
		}
	}
	return cell(column, value, byteColumns)
}

func compactJSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(encoded)
}
