package pretty

import (
	"fmt"
	"io"
	"sort"
	"strings"

	prettytable "github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const maxCellWidth = 48

func renderRows(w io.Writer, tableTitle string, rows []map[string]any) error {
	byteColumns := inferByteColumns(rows)
	columns := collectColumns(rows)
	constants, columns := extractConstants(rows, columns)
	sortColumns(rows, columns)
	if len(columns) == 0 {
		return writeMetadata(w, constants, byteColumns)
	}

	headings := make(prettytable.Row, len(columns))
	for i, column := range columns {
		headings[i] = header(column)
	}
	tableRows := make([]prettytable.Row, len(rows))
	for i, row := range rows {
		tableRows[i] = make(prettytable.Row, len(columns))
		for j, column := range columns {
			tableRows[i][j] = tableCell(column, row[column], byteColumns)
		}
	}

	tw := prettytable.NewWriter()
	tw.SetStyle(prettytable.StyleRounded)
	tw.Style().Format.Header = text.FormatDefault
	tw.Style().Options.SeparateRows = true
	tw.AppendHeader(headings)
	tw.AppendRows(tableRows)
	if tableTitle != "" {
		tw.SetTitle(tableTitle)
	}
	columnConfigs := make([]prettytable.ColumnConfig, len(columns))
	for i, column := range columns {
		columnConfigs[i] = prettytable.ColumnConfig{
			Align:            numericColumnAlignment(rows, column),
			Number:           i + 1,
			WidthMax:         maxCellWidth,
			WidthMaxEnforcer: text.WrapSoft,
		}
	}
	tw.SetColumnConfigs(columnConfigs)
	rendered := tw.Render()
	if len(constants) > 0 {
		keys := make([]string, 0, len(constants))
		for key := range constants {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var details strings.Builder
		for _, key := range keys {
			fmt.Fprintf(&details, "├── %s: %s\n", title(key), cell(key, constants[key], byteColumns))
		}
		details.WriteString("│\n├")
		details.WriteString(strings.TrimPrefix(rendered, "╭"))
		rendered = details.String()
	}
	_, err := fmt.Fprintln(w, rendered)
	return err
}

func numericColumnAlignment(rows []map[string]any, column string) text.Align {
	seen := false
	for _, row := range rows {
		value, exists := row[column]
		if !exists || value == nil {
			continue
		}
		if _, ok := number(value); !ok {
			return text.AlignDefault
		}
		seen = true
	}
	if seen {
		return text.AlignRight
	}
	return text.AlignDefault
}
