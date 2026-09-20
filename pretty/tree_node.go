package pretty

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

func renderTreeNode(out *strings.Builder, name string, value any, prefix string, last bool) error {
	connector, childPrefix := "├── ", prefix+"│   "
	if last {
		connector, childPrefix = "└── ", prefix+"    "
	}

	switch data := value.(type) {
	case map[string]any:
		if len(data) == 0 {
			fmt.Fprintf(out, "%s%s%s: {}\n", prefix, connector, title(name))
			return nil
		}
		fmt.Fprintf(out, "%s%s%s\n", prefix, connector, title(name))
		return renderTreeObject(out, data, childPrefix)
	case []any:
		if rows, ok := recordCollection(data); ok {
			if len(rows) == 0 {
				fmt.Fprintf(out, "%s%s%s: []\n", prefix, connector, title(name))
				return nil
			}
			if !rowsAreFlat(rows) {
				fmt.Fprintf(out, "%s%s%s\n", prefix, connector, title(name))
				for i, row := range rows {
					if err := renderTreeNode(out, fmt.Sprintf("Item %d", i+1), row, childPrefix, i == len(rows)-1); err != nil {
						return err
					}
				}
				return nil
			}

			var table strings.Builder
			if err := renderRows(&table, title(name), rows); err != nil {
				return err
			}
			rendered := strings.TrimSuffix(table.String(), "\n")
			lines := strings.Split(rendered, "\n")
			if strings.HasPrefix(lines[0], "╭") {
				lines[0] = "├" + strings.TrimPrefix(lines[0], "╭")
			}
			if !last && strings.HasPrefix(lines[len(lines)-1], "╰") {
				lines[len(lines)-1] = "├" + strings.TrimPrefix(lines[len(lines)-1], "╰")
			}
			for _, line := range lines {
				fmt.Fprintf(out, "%s%s\n", prefix, line)
			}
			return nil
		}
	}

	byteColumns := inferByteColumns([]map[string]any{{name: value}})
	lines := strings.Split(text.WrapSoft(cell(name, value, byteColumns), maxCellWidth), "\n")
	fmt.Fprintf(out, "%s%s%s: %s\n", prefix, connector, title(name), lines[0])
	for _, line := range lines[1:] {
		fmt.Fprintf(out, "%s    %s\n", childPrefix, line)
	}
	return nil
}
