package pretty

import (
	"sort"
	"strings"
)

func renderTreeObject(out *strings.Builder, object map[string]any, prefix string) error {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		leftLeaf, rightLeaf := treeLeaf(object[keys[i]]), treeLeaf(object[keys[j]])
		if leftLeaf != rightLeaf {
			return leftLeaf
		}
		return keys[i] < keys[j]
	})
	for i, key := range keys {
		if err := renderTreeNode(out, key, object[key], prefix, i == len(keys)-1); err != nil {
			return err
		}
	}
	return nil
}

func treeLeaf(value any) bool {
	switch data := value.(type) {
	case map[string]any:
		return false
	case []any:
		_, records := recordCollection(data)
		return !records
	default:
		return true
	}
}
