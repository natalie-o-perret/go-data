package data

// Transform clones map[string]any and []any containers and calls transform on
// each value after transforming its descendants. Map keys are visited in
// lexical order and slice elements in index order.
func Transform(value any, transform TransformFunc) (any, error) {
	active := make(map[containerID]bool)

	var apply func(Path, any) (any, error)
	apply = func(path Path, value any) (any, error) {
		if id, container := identify(value); container {
			if active[id] {
				return nil, ErrRecursive
			}
			active[id] = true
			defer delete(active, id)
		}

		transformed := value
		switch value := value.(type) {
		case map[string]any:
			if value != nil {
				clone := make(map[string]any, len(value))
				for _, key := range sortedKeys(value) {
					child, err := apply(childPath(path, key), value[key])
					if err != nil {
						return nil, err
					}
					clone[key] = child
				}
				transformed = clone
			}
		case []any:
			if value != nil {
				clone := make([]any, len(value))
				for index, child := range value {
					result, err := apply(childPath(path, index), child)
					if err != nil {
						return nil, err
					}
					clone[index] = result
				}
				transformed = clone
			}
		}
		return transform(path, transformed)
	}

	return apply(nil, value)
}
