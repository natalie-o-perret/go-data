// Package pretty renders decoded structured data as deterministic trees and tables.
package pretty

import (
	"errors"
	"fmt"
	"io"

	data "github.com/natalie-o-perret/go-data"
	"github.com/natalie-o-perret/go-data/json"
)

// Render writes value as a deterministic tree or table, falling back to JSON.
// Recursive maps and slices return an error before any output is written.
func Render(w io.Writer, value any) error {
	value = normalizeRoot(value)
	if err := data.Walk(value, func(data.Path, any) error { return nil }); err != nil {
		if errors.Is(err, data.ErrRecursive) {
			return fmt.Errorf("pretty: recursive data structures are not supported: %w", err)
		}
		return err
	}
	if requiresTree(value) {
		return renderTree(w, value)
	}
	metadata, sections, ok := tableSections(value)
	if !ok {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(value)
	}
	if len(metadata) > 0 {
		if err := writeMetadata(w, metadata, inferByteColumns([]map[string]any{metadata})); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	for i, section := range sections {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if len(section.rows) == 0 {
			if section.title != "" {
				if _, err := fmt.Fprintln(w, title(section.title)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w, "(empty)"); err != nil {
				return err
			}
			continue
		}
		if err := renderRows(w, title(section.title), section.rows); err != nil {
			return err
		}
	}
	return nil
}

func normalizeRoot(value any) any {
	if items, ok := value.([]any); ok && len(items) == 1 {
		if _, ok := items[0].(map[string]any); ok {
			return items[0]
		}
	}
	return value
}
