package pretty

import "os"

func ExampleRender() {
	value := map[string]any{
		"catalog": map[string]any{
			"name": "Example Press",
			"owner": map[string]any{
				"name": "Alex",
				"contact": map[string]any{
					"email": "alex@example.test",
				},
			},
			"sections": []any{
				map[string]any{
					"name": "Fiction",
					"publications": []any{
						map[string]any{
							"id": "pub-001", "title": "Field Notes",
							"file_bytes": 8 << 20,
							"formats":    []any{"print", "ebook"},
						},
						map[string]any{
							"id": "pub-002", "title": "Short Stories",
							"file_bytes": 4 << 20,
							"formats":    []any{"ebook"},
						},
					},
					"awards": []any{
						map[string]any{"name": "Readers' Choice", "year": 2026},
					},
				},
				map[string]any{
					"name": "Essays",
					"publications": []any{
						map[string]any{
							"id": "pub-003", "title": "Small Observations",
							"file_bytes": 2 << 20,
							"formats":    []any{"print"},
						},
					},
				},
			},
		},
		"warnings": []any{
			map[string]any{"code": "stale", "message": "One source is delayed"},
			map[string]any{"code": "draft", "message": "One publication is unpublished"},
		},
	}
	if err := Render(os.Stdout, value); err != nil {
		panic(err)
	}

	// Output:
	// ├── Catalog
	// │   ├── Name: Example Press
	// │   ├── Owner
	// │   │   ├── Name: Alex
	// │   │   └── Contact
	// │   │       └── Email: alex@example.test
	// │   └── Sections
	// │       ├── Item 1
	// │       │   ├── Name: Fiction
	// │       │   ├────────────────────────╮
	// │       │   │ Awards                 │
	// │       │   ├─────────────────┬──────┤
	// │       │   │ Name            │ Year │
	// │       │   ├─────────────────┼──────┤
	// │       │   │ Readers' Choice │ 2026 │
	// │       │   ├─────────────────┴──────╯
	// │       │   ├─────────────────────────────────────────────────────╮
	// │       │   │ Publications                                        │
	// │       │   ├────────────┬─────────┬───────────────┬──────────────┤
	// │       │   │ File Bytes │ Id      │ Title         │ Formats      │
	// │       │   ├────────────┼─────────┼───────────────┼──────────────┤
	// │       │   │      8 MiB │ pub-001 │ Field Notes   │ print, ebook │
	// │       │   ├────────────┼─────────┼───────────────┼──────────────┤
	// │       │   │      4 MiB │ pub-002 │ Short Stories │ ebook        │
	// │       │   ╰────────────┴─────────┴───────────────┴──────────────╯
	// │       └── Item 2
	// │           ├── Name: Essays
	// │           ├─────────────────────────────────────────────────────╮
	// │           │ Publications                                        │
	// │           ├────────────┬─────────┬────────────────────┬─────────┤
	// │           │ File Bytes │ Id      │ Title              │ Formats │
	// │           ├────────────┼─────────┼────────────────────┼─────────┤
	// │           │      2 MiB │ pub-003 │ Small Observations │ print   │
	// │           ╰────────────┴─────────┴────────────────────┴─────────╯
	// ├────────────────────────────────────────╮
	// │ Warnings                               │
	// ├───────┬────────────────────────────────┤
	// │ Code  │ Message                        │
	// ├───────┼────────────────────────────────┤
	// │ stale │ One source is delayed          │
	// ├───────┼────────────────────────────────┤
	// │ draft │ One publication is unpublished │
	// ╰───────┴────────────────────────────────╯
}
