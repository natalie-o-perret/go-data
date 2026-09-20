package data

// Walk visits value and then its map[string]any and []any descendants.
// Map keys are visited in lexical order and slice elements in index order.
func Walk(value any, visit VisitFunc) error {
	active := make(map[containerID]bool)

	var walk func(Path, any) error
	walk = func(path Path, value any) error {
		if id, container := identify(value); container {
			if active[id] {
				return ErrRecursive
			}
			active[id] = true
			defer delete(active, id)
		}

		if err := visit(path, value); err != nil {
			return err
		}

		switch value := value.(type) {
		case map[string]any:
			for _, key := range sortedKeys(value) {
				if err := walk(childPath(path, key), value[key]); err != nil {
					return err
				}
			}
		case []any:
			for index, child := range value {
				if err := walk(childPath(path, index), child); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return walk(nil, value)
}
