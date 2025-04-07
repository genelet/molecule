package godbi

// cloneArgs clones args to a new args, keeping proper data type
func cloneArgs(args any) any {
	if args == nil {
		return nil
	}
	switch t := args.(type) {
	case map[string]any:
		return cloneMap(t)
	case []map[string]any:
		var newArgs []any
		for _, each := range t {
			newArgs = append(newArgs, each)
		}
		return newArgs
	case []any:
		var newArgs []any
		for _, each := range t {
			if item, ok := each.(map[string]any); ok {
				newArgs = append(newArgs, item)
			}
		}
		return newArgs
	default:
	}
	return nil
}

func mergeArgsMap(args any, item map[string]any, force ...bool) any {
	if args == nil {
		return item
	} else if item == nil {
		return args
	}
	switch t := args.(type) {
	case map[string]any:
		if force != nil && force[0] {
			return mergeMap(t, item)
		}
		return mergeMapOr(t, item)
	case []map[string]any:
		var newArgs []any
		for _, each := range t {
			if force != nil && force[0] {
				newArgs = append(newArgs, mergeMap(each, item))
			} else {
				newArgs = append(newArgs, mergeMapOr(each, item))
			}
		}
		return newArgs
	case []any:
		var newArgs []any
		for _, each := range t {
			if single, ok := each.(map[string]any); ok {
				if force != nil && force[0] {
					newArgs = append(newArgs, mergeMap(single, item))
				} else {
					newArgs = append(newArgs, mergeMapOr(single, item))
				}
			}
		}
		return newArgs
	default:
	}
	return nil
}

// mergeArgs merges map to either an existing map, or slice of map in which each element will be merged
func mergeArgs(args, items any, force ...bool) any {
	if !hasValue(args) {
		return items
	} else if !hasValue(items) {
		return args
	}

	middle := func(newArgs *[]any, args any, item map[string]any, force ...bool) {
		merged := mergeArgsMap(args, item, force...)
		switch s := merged.(type) {
		case []any:
			for _, u := range s {
				*newArgs = append(*newArgs, u.(map[string]any))
			}
		case []map[string]any:
			for _, u := range s {
				*newArgs = append(*newArgs, u)
			}
		default:
			*newArgs = append(*newArgs, merged.(map[string]any))
		}
	}

	switch t := items.(type) {
	case map[string]any:
		return mergeArgsMap(args, t, force...)
	case []map[string]any:
		var newArgs = make([]any, 0)
		for _, item := range t {
			middle(&newArgs, args, item, force...)
		}
		return newArgs
	case []any:
		var newArgs = make([]any, 0)
		for _, each := range t {
			item, ok := each.(map[string]any)
			if !ok {
				continue
			}
			middle(&newArgs, args, item, force...)
		}
		return newArgs
	default:
	}
	return args
}
