package godbi

// cloneArgs clones args to a new args, keeping proper data type
func cloneArgs(args interface{}) interface{} {
	if args == nil {
		return nil
	}
	switch t := args.(type) {
	case map[string]interface{}:
		return cloneMap(t)
	case []map[string]interface{}:
		var newArgs []interface{}
		for _, each := range t {
			newArgs = append(newArgs, each)
		}
		return newArgs
	case []interface{}:
		var newArgs []interface{}
		for _, each := range t {
			if item, ok := each.(map[string]interface{}); ok {
				newArgs = append(newArgs, item)
			}
		}
		return newArgs
	default:
	}
	return nil
}

func mergeArgsMap(args interface{}, item map[string]interface{}, force ...bool) interface{} {
	if args == nil {
		return item
	} else if item == nil {
		return args
	}
	switch t := args.(type) {
	case map[string]interface{}:
		if force != nil && force[0] {
			return mergeMap(t, item)
		}
		return mergeMapOr(t, item)
	case []map[string]interface{}:
		var newArgs []interface{}
		for _, each := range t {
			if force != nil && force[0] {
				newArgs = append(newArgs, mergeMap(each, item))
			} else {
				newArgs = append(newArgs, mergeMapOr(each, item))
			}
		}
		return newArgs
	case []interface{}:
		var newArgs []interface{}
		for _, each := range t {
			if single, ok := each.(map[string]interface{}); ok {
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
func mergeArgs(args, items interface{}, force ...bool) interface{} {
	if !hasValue(args) {
		return items
	} else if !hasValue(items) {
		return args
	}

	middle := func(newArgs *[]interface{}, args interface{}, item map[string]interface{}, force ...bool) {
		merged := mergeArgsMap(args, item, force...)
		switch s := merged.(type) {
		case []interface{}:
			for _, u := range s {
				*newArgs = append(*newArgs, u.(map[string]interface{}))
			}
		case []map[string]interface{}:
			for _, u := range s {
				*newArgs = append(*newArgs, u)
			}
		default:
			*newArgs = append(*newArgs, merged.(map[string]interface{}))
		}
	}

	switch t := items.(type) {
	case map[string]interface{}:
		return mergeArgsMap(args, t, force...)
	case []map[string]interface{}:
		var newArgs = make([]interface{}, 0)
		for _, item := range t {
			middle(&newArgs, args, item, force...)
		}
		return newArgs
	case []interface{}:
		var newArgs = make([]interface{}, 0)
		for _, each := range t {
			item, ok := each.(map[string]interface{})
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
