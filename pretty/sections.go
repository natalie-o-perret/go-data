package pretty

import "sort"

type tableSection struct {
	title string
	rows  []map[string]any
}

func tableSections(value any) (map[string]any, []tableSection, bool) {
	if object, ok := value.(map[string]any); ok {
		metadata := map[string]any{}
		var sections []tableSection
		for title, candidate := range object {
			if rows, ok := recordCollection(candidate); ok {
				sections = append(sections, tableSection{title: title, rows: rows})
			} else {
				metadata[title] = candidate
			}
		}
		if len(sections) > 0 {
			sort.Slice(sections, func(i, j int) bool { return sections[i].title < sections[j].title })
			return metadata, sections, true
		}
	}
	rows, ok := tableRows(value)
	if !ok {
		return nil, nil, false
	}
	return nil, []tableSection{{rows: rows}}, true
}

func recordCollection(value any) ([]map[string]any, bool) {
	items, ok := value.([]any)
	if !ok {
		return nil, false
	}
	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		row, ok := item.(map[string]any)
		if !ok {
			return nil, false
		}
		rows = append(rows, row)
	}
	return rows, true
}

func tableRows(value any) ([]map[string]any, bool) {
	switch data := value.(type) {
	case map[string]any:
		return []map[string]any{data}, true
	case []any:
		return recordCollection(data)
	default:
		return nil, false
	}
}
