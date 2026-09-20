package pretty

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

func number(value any) (float64, bool) {
	switch n := value.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		parsed, err := n.Float64()
		return parsed, err == nil
	default:
		return 0, false
	}
}

func inferByteColumns(rows []map[string]any) map[string]bool {
	result := map[string]bool{}
	for _, column := range collectColumns(rows) {
		normalized := strings.ToLower(strings.ReplaceAll(column, "-", "_"))
		if normalized == "bytes" || strings.HasSuffix(normalized, "_bytes") {
			result[column] = true
			continue
		}
		if normalized != "memory" && !strings.HasSuffix(normalized, "_memory") {
			continue
		}
		seen := false
		valid := true
		for _, row := range rows {
			value, exists := row[column]
			if !exists || value == nil {
				continue
			}
			n, ok := number(value)
			if !ok || math.Abs(n) < 1<<20 || math.Mod(math.Abs(n), 1024) != 0 {
				valid = false
				break
			}
			seen = true
		}
		if seen && valid {
			result[column] = true
		}
	}
	return result
}

func humanBytes(bytes float64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	value := bytes
	unit := 0
	for math.Abs(value) >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	return strconv.FormatFloat(value, 'f', -1, 64) + " " + units[unit]
}
