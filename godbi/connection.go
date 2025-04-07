package godbi

type ConnectType int

const (
	CONNECTDefault ConnectType = iota
	CONNECTOne
	CONNECTArray
	CONNECTMap
	CONNECTMany
)

// Connection describes linked page
// 1) for Nextpages, it maps item in lists to next ARGS and next Extra
// 2) for Prepares, it maps current ARGS to the next ARGS and next Extra
type Connection struct {
	// AtomName: the name of the table
	AtomName string `json:"atomName" hcl:"atomName,label"`

	// ActionName: the action on the atom
	ActionName string `json:"actionName" hcl:"actionName,label"`

	// RelateArgs: map current page's columns to nextpage's columns as input
	RelateArgs map[string]string `json:"relateArgs,omitempty" hcl:"relateArgs,optional"`

	// RelateExtra: map current page's columns to nextpage's columns (for Nextpages), or prepared page's columns to current page's columns (for Prepares) as constrains.
	RelateExtra map[string]string `json:"relateExtra,omitempty" hcl:"relateExtra,optional"`

	// Dimension: for nextpage's output format
	Dimension ConnectType `json:"dimension,omitempty" hcl:"dimension,optional"`

	// Marker: for input data, this marks a whole data set for the next or previous object; for output data, this is the key for the next whole data set.
	Marker string `json:"marker,omitempty" hcl:"marker,optional"`
}

// Subname is the marker string used to store the output
func (self *Connection) Subname() string {
	if self.Marker != "" {
		return self.Marker
	}
	return self.AtomName + "_" + self.ActionName
}

// findExrea returns the value if the input i.e. item contains
// the current table name as key.
func (self *Connection) findExrea(item map[string]any) map[string]any {
	marker := self.Marker
	if marker == "" {
		return nil
	}

	if v, ok := item[marker]; ok {
		switch t := v.(type) {
		case map[string]any:
			return t
		default:
		}
	}
	return nil
}

// findArgs returns the value if the input i.e. args contains
// the current table name as key.
func (self *Connection) findArgs(args any) (any, bool) {
	if args == nil {
		return nil, true
	}

	marker := self.Marker
	if marker == "" {
		return nil, false
	}

	switch t := args.(type) {
	case map[string]any: // in practice, only this data type exists
		if v, ok := t[marker]; ok {
			switch s := v.(type) {
			case map[string]any:
				if self.Dimension == CONNECTMap {
					var outs []map[string]any
					for key, value := range s {
						outs = append(outs, map[string]any{"key": key, "value": value})
					}
					return outs, true
				}
				return s, true
			case []map[string]any:
				return s, true
			case []any:
				var outs []map[string]any
				for _, item := range s {
					switch x := item.(type) {
					case map[string]any:
						outs = append(outs, x)
					default: // native types
						outs = append(outs, map[string]any{marker: x})
					}
				}
				return outs, true
			default:
			}
			return nil, true
		}
	default:
	}

	return nil, true
}

// NextArg returns nextpage's args as the value of key using current args map
func (self *Connection) nextArgs(args any) any {
	if args == nil {
		return nil
	}
	if _, ok := self.RelateArgs["ALL"]; ok {
		smallerRelate := make(map[string]string)
		for k, v := range self.RelateArgs {
			if k == "ALL" {
				continue
			}
			smallerRelate[k] = v
		}
		if len(smallerRelate) == 0 {
			return args
		}
		smallerArgs := nextArgsFromRelate(args, smallerRelate)
		// keys in args will overrde exising keys in smallerArgs
		return mergeArgs(smallerArgs, args, true)
	}
	return nextArgsFromRelate(args, self.RelateArgs)
}

func nextArgsFromRelate(args any, relate map[string]string) any {
	if args == nil {
		return nil
	}
	switch t := args.(type) {
	case map[string]any:
		return createNextmap(relate, t)
	case []map[string]any:
		var outs []any
		for _, hash := range t {
			if x := createNextmap(relate, hash); x != nil {
				outs = append(outs, x)
			}
		}
		return outs
	case []any:
		var outs []any
		for _, hash := range t {
			if item, ok := hash.(map[string]any); ok {
				if x := createNextmap(relate, item); x != nil {
					outs = append(outs, x)
				}
			}
		}
		return outs
	default:
	}
	return nil
}

// nextExtra returns nextpage's extra using current extra map
func (self *Connection) nextExtra(args any) map[string]any {
	if _, ok := self.RelateExtra["ALL"]; ok {
		if v, ok := args.(map[string]any); ok {
			return v
		}
		return nil
	}

	switch t := args.(type) {
	case map[string]any:
		return createNextmap(self.RelateExtra, t)
	default:
	}
	return nil
}

func createNextmap(which map[string]string, item map[string]any) map[string]any {
	if which == nil {
		return nil
	}

	var args map[string]any
	for k, v := range which {
		if u, ok := item[k]; ok {
			if args == nil {
				args = make(map[string]any)
			}
			switch t := u.(type) {
			case map[string]any:
				for key, value := range t {
					args[key] = value
				}
			default:
				args[v] = t
			}
		}
	}
	return args
}

func (self *Connection) shorten(lists []any) any {
	if self.Dimension == CONNECTDefault || self.Marker == "" {
		return lists
	}

	switch self.Dimension {
	case CONNECTMap:
		output := make(map[string]any)
		for _, single := range lists {
			key, value := mapEntryToPair(single.(map[string]any))
			if key != "" {
				output[key] = value
			}
		}
		return output
	case CONNECTMany:
		var output []any
		for _, single := range lists {
			output = append(output, manyEntry(self.Marker, single.(map[string]any)))
		}
		return output
	case CONNECTArray:
		var output []any
		for _, single := range lists {
			output = append(output, single.(map[string]any)[self.Marker])
		}
		return output
	case CONNECTOne:
		return manyEntry(self.Marker, lists[0].(map[string]any))
	default:
	}
	return lists
}

func mapEntryToPair(single map[string]any) (string, any) {
	var key string
	var value any
	var object map[string]any

	for _, a := range single {
		var c map[string]any
		switch b := a.(type) {
		case []any:
			c = b[0].(map[string]any)
		case []map[string]any:
			c = b[0]
		default:
			continue
		}

		for d, e := range c {
			if d == "key" {
				key = e.(string)
				continue
			}
			if d == "value" {
				value = e
				continue
			}
			switch f := e.(type) {
			case []any:
				object = f[0].(map[string]any)
			case []map[string]any:
				object = f[0]
			default:
			}
		}
	}

	if value != nil {
		return key, value
	}
	return key, object
}

func manyEntry(leader string, single map[string]any) map[string]any {
	output := make(map[string]any)
	var higher map[string]any

	for key, value := range single {
		if key == leader {
			higher = make(map[string]any)
			switch t := value.(type) {
			case []any:
				for k, v := range t[0].(map[string]any) {
					higher[k] = v
				}
			case []map[string]any:
				for k, v := range t[0] {
					higher[k] = v
				}
			case map[string]any:
				for k, v := range t {
					higher[k] = v
				}
			default:
				higher[key] = value
			}
		} else {
			output[key] = value
		}
	}

	for k, v := range higher {
		output[k] = v
	}

	return output
}

func shortRecursive(leader string, single map[string]any) map[string]any {
	output := make(map[string]any)
	var higher map[string]any

	for key, value := range single {
		if key == leader {
			higher = make(map[string]any)
			switch t := value.(type) {
			case []any:
				if len(t) > 1 {
					higher[key] = t
					continue
				} else {
					for _, s := range t {
						for k, v := range s.(map[string]any) {
							switch u := v.(type) {
							case []any:
								old, ok := higher[k]
								if ok {
									higher[k] = append(old.([]any), u...)
								} else {
									higher[k] = u
								}
							default:
								higher[k] = v
							}
						}
					}
				}
			case []map[string]any:
				for k, v := range t[0] {
					higher[k] = v
				}
			case map[string]any:
				for k, v := range t {
					higher[k] = v
				}
			default:
				higher[key] = value
			}
		} else {
			output[key] = value
		}
	}

	for k, v := range higher {
		output[k] = v
	}

	return output
}

func (self *Connection) shortenRecursive(lists []any) any {
	switch self.Dimension {
	case CONNECTOne:
		return shortRecursive(self.Marker, lists[0].(map[string]any))
	default:
	}

	var output []any
	for _, single := range lists {
		output = append(output, shortRecursive(self.Marker, single.(map[string]any)))
	}
	return output
}

// useless but put here for backup
func shortenX(leader string, lists []any) []any {
	extra := make(map[string]any)

	var higher map[string]any
	var temp []any

	for key, value := range lists[0].(map[string]any) {
		if key == leader {
			higher = make(map[string]any)
			switch t := value.(type) {
			case []any:
				if len(t) > 1 {
					higher[key] = t
				} else {
					for _, s := range t {
						tmp := make(map[string]any)
						for k, v := range s.(map[string]any) {
							switch u := v.(type) {
							case []any:
								old, ok := tmp[k]
								if ok {
									tmp[k] = append(old.([]any), u...)
								} else {
									tmp[k] = u
								}
							default:
								tmp[k] = v
							}
						}
						temp = append(temp, tmp)
					}
				}
			case []map[string]any:
				higher = make(map[string]any)
				for k, v := range t[0] {
					higher[k] = v
				}
			case map[string]any:
				higher = make(map[string]any)
				for k, v := range t {
					higher[k] = v
				}
			default:
				higher = make(map[string]any)
				higher[key] = value
			}
		} else {
			extra[key] = value
		}
	}

	if hasValue(temp) {
		var output []any
		for _, tmp := range temp {
			for k, v := range tmp.(map[string]any) {
				extra[k] = v
			}
			output = append(output, extra)
		}
		return output
	}

	for k, v := range higher {
		extra[k] = v
	}
	return []any{extra}
}
