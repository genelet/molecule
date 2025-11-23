package godbi

import (
	"context"
	"database/sql"
)

type PreStopper interface {
	Stop(tableObj, childObj *Table) bool
}

type Stopper interface {
	PreStopper
	Sign(tableObj *Table, item map[string]any)
}

// Molecule describes all atoms and actions in a database schema
type Molecule struct {
	Atoms    []*Atom `json:"atoms" hcl:"atoms,block"`
	DBDriver DBType  `json:"dbDriver" hcl:"dbDriver,optional"`
	Stopper
	PreStopper
	logger Slogger
}

// SetLogger sets the logger
func (m *Molecule) SetLogger(logger Slogger) {
	m.logger = logger
	for _, atom := range m.Atoms {
		atom.Table.SetLogger(logger)
	}
}

// GetLogger gets the logger
func (m *Molecule) GetLogger() Slogger {
	return m.logger
}

// GetAtom returns the atom by atom name
func (m *Molecule) GetAtom(atomName string) *Atom {
	if m.Atoms != nil {
		for _, atom := range m.Atoms {
			if atom.AtomName == atomName {
				atom.SetDBDriver(m.DBDriver)
				atom.Table.logger = m.logger
				return atom
			}
		}
	}
	return nil
}

// RunOption holds the arguments for RunContext
type RunOption struct {
	Args        any
	Extra       map[string]any
	GlobalArgs  map[string]any
	GlobalExtra map[string]any
}

// RunContext runs action by atom and action string names.
// It returns the searched data and optional error code.
// atom is the atom name, and action the action name.
func (m *Molecule) RunContext(ctx context.Context, db *sql.DB, atom, action string, opt *RunOption) ([]any, error) {
	return m.processContext(false, ctx, db, atom, action, opt)
}

func (m *Molecule) runRecurseContext(ctx context.Context, db *sql.DB, atom, action string, opt *RunOption) ([]any, error) {
	return m.processContext(true, ctx, db, atom, action, opt)
}

func (m *Molecule) processContext(topRecursive bool, ctx context.Context, db *sql.DB, atom, action string, opt *RunOption) ([]any, error) {
	var args any
	var extra map[string]any
	var globalArgs map[string]any
	var globalExtra map[string]any

	if opt != nil {
		args = opt.Args
		extra = opt.Extra
		globalArgs = opt.GlobalArgs
		globalExtra = opt.GlobalExtra
	}

	if hasValue(globalArgs) && hasValue(globalArgs[atom]) {
		argsMap := globalArgs[atom].(map[string]any)
		args = mergeArgs(args, argsMap[action])
	}

	if hasValue(globalExtra) && hasValue(globalExtra[atom]) {
		extraAction := globalExtra[atom].(map[string]any)
		if hasValue(extraAction[action]) {
			extra = mergeMap(extra, extraAction[action].(map[string]any))
		}
	}

	switch t := args.(type) {
	case map[string]any:
		return m.execContext(topRecursive, ctx, db, atom, action, t, extra, globalArgs, globalExtra)
	case []map[string]any:
		var final []any
		for _, arg := range t {
			lists, err := m.execContext(topRecursive, ctx, db, atom, action, arg, extra, globalArgs, globalExtra)
			if err != nil {
				return nil, err
			}
			final = append(final, lists...)
		}
		return final, nil
	case []any:
		var final []any
		for _, arg := range t {
			if v, ok := arg.(map[string]any); ok {
				lists, err := m.execContext(topRecursive, ctx, db, atom, action, v, extra, globalArgs, globalExtra)
				if err != nil {
					return nil, err
				}
				final = append(final, lists...)
			}
		}
		return final, nil
	default:
	}

	return m.execContext(topRecursive, ctx, db, atom, action, nil, extra, globalArgs, globalExtra)
}

// execContext executes the action logic with fully resolved arguments.
func (m *Molecule) execContext(topRecursive bool, ctx context.Context, db *sql.DB, atom, action string, args, extra, globalArgs, globalExtra map[string]any) ([]any, error) {
	atomObj := m.GetAtom(atom)
	if atomObj == nil {
		return nil, errorAtomNotFound(atom)
	}
	tableObj := atomObj.Table
	actionObj := atomObj.GetAction(action)
	if actionObj == nil {
		return nil, errorActionNotFound(action, atom)
	}
	if actionObj.GetBaseAction().IsDo && !hasValue(args) {
		return nil, nil
	}

	prepares := actionObj.GetBaseAction().Prepares
	nextpages := actionObj.GetBaseAction().Nextpages

	newArgs := cloneArgs(args)
	newExtra := cloneMap(extra)

	// prepares receives filtered args and extra from current args
	for _, p := range prepares {
		pAtom := m.GetAtom(p.AtomName)
		if pAtom == nil {
			return nil, errorAtomNotFound(p.AtomName)
		}
		pTable := pAtom.Table
		v, _ := p.findArgs(args)
		preArgs := mergeArgs(p.nextArgs(args), v)
		preExtra := mergeMap(p.nextExtra(args), p.findExrea(extra))
		isDo := pAtom.GetAction(p.ActionName).GetBaseAction().IsDo
		isRecursive := pTable.IsRecursive()

		var lists []any
		var err error
		if topRecursive {
			if !hasValue(preArgs) {
				return []any{args}, nil
			}
			pk := pTable.Pks[0]
			if isRecursive {
				switch t := preArgs.(type) {
				case map[string]any:
					delete(t, pk)
				case []any:
					for _, s := range t {
						delete(s.(map[string]any), pk)
					}
				case []map[string]any:
					for _, s := range t {
						delete(s, pk)
					}
				}
			}
			lists, err = m.runRecurseContext(ctx, db, p.AtomName, p.ActionName, &RunOption{Args: preArgs, Extra: preExtra, GlobalArgs: globalArgs, GlobalExtra: globalExtra})
		} else if isDo && isRecursive {
			// this triggers the original topRecursive and is always a DO action
			lists, err = m.runRecurseContext(ctx, db, p.AtomName, p.ActionName, &RunOption{Args: preArgs, Extra: preExtra, GlobalArgs: globalArgs, GlobalExtra: globalExtra})
		} else if m.PreStopper == nil || !m.PreStopper.Stop(&tableObj, &pTable) {
			lists, err = m.RunContext(ctx, db, p.AtomName, p.ActionName, &RunOption{Args: preArgs, Extra: preExtra, GlobalArgs: globalArgs, GlobalExtra: globalExtra})
		} // else means stopper is set and stopper stops preparing insert or update
		if err != nil {
			return nil, err
		}

		// only two types of prepares
		// 1) one pre, with multiple outputs (when p.argsMap is multiple)
		if hasValue(lists) && len(lists) > 1 {
			var tmp []map[string]any
			newExtra = cloneMap(extra)
			for _, item := range lists {
				result := mergeArgs(args, p.nextArgs(item)).(map[string]any)
				tmp = append(tmp, result)
				newExtra = mergeMap(newExtra, p.nextExtra(item))
			}
			newArgs = tmp
			if p.ActionName == "delecs" {
				continue
			} else if isDo {
				break
			}
		}
		// 2) multiple pre, with one output each.
		// when a multiple output is found, 1) will override
		if hasValue(lists) && hasValue(lists[0]) {
			middle := p.nextArgs(lists[0])
			if topRecursive && isRecursive {
				pk := pTable.Pks[0]
				switch t := middle.(type) {
				case map[string]any:
					delete(t, pk)
				case []any:
					for _, s := range t {
						delete(s.(map[string]any), pk)
					}
				case []map[string]any:
					for _, s := range t {
						delete(s, pk)
					}
				}
			}
			newArgs = mergeArgs(newArgs, middle.(map[string]any), true)
			newExtra = mergeMap(newExtra, p.nextExtra(lists[0]))
		}
	}

	if !topRecursive && hasValue(newArgs) && actionObj.GetBaseAction().IsDo {
		//newArgs = tableObj.refreshArgs(newArgs, true)
		newArgs = tableObj.refreshArgs(newArgs)
	}

	data, err := atomObj.RunAtomContext(ctx, db, action, newArgs, newExtra)
	if err != nil {
		return nil, err
	}

	if actionObject, ok := actionObj.(*Topics); ok {
		switch t := newArgs.(type) {
		case map[string]any:
			for _, item := range []string{actionObject.FIELDS, actionObject.SORTBY, actionObject.SORTREVERSE, actionObject.PAGESIZE, actionObject.PAGENO, actionObject.TOTALNO, actionObject.MAXPAGENO} {
				if _, ok := t[item]; ok {
					if _, ok := args[item]; !ok {
						args[item] = t[item]
					}
				}
			}
		default:
		}
	}

	for _, item := range data {
		if item == nil {
			continue
		}
		hash, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if m.Stopper != nil {
			m.Stopper.Sign(&tableObj, hash)
		}
	}

	if topRecursive && action != "delecs" {
		if tableObj.IsRecursive() {
			var p *Connection
			var pAtom *Atom
			var rColumn string
			for _, p = range nextpages {
				pAtom = m.GetAtom(p.AtomName)
				if pAtom == nil {
					return nil, errorAtomNotFound(p.AtomName)
				}
				rColumn = pAtom.Table.RecursiveColumn()
				if rColumn != "" {
					break
				}
			}
			if rColumn == "" {
				return data, nil
			}
			argsData, _ := p.findArgs(newArgs)
			if !hasValue(argsData) {
				return data, nil
			}
			for _, item := range data {
				nextArgs := p.nextArgs(item)
				nextArgs = mergeArgs(argsData, nextArgs)
				_, err = m.runRecurseContext(ctx, db, p.AtomName, p.ActionName, &RunOption{Args: nextArgs, GlobalArgs: globalArgs, GlobalExtra: globalExtra})
				if err != nil {
					return nil, err
				}
			}
		}
		return data, err
	}

	for _, p := range nextpages {
		pAtom := m.GetAtom(p.AtomName)
		if pAtom == nil {
			return nil, errorAtomNotFound(p.AtomName)
		}
		pAction := pAtom.GetAction(p.ActionName)
		if m.Stopper != nil && m.Stopper.Stop(&tableObj, &(pAtom.Table)) {
			continue
		}

		for _, item := range data {
			if item == nil {
				continue
			}
			nextArgs := p.nextArgs(item)
			nextExtra := p.nextExtra(item)
			if pAction.GetBaseAction().IsDo {
				if v, ok := p.findArgs(newArgs); ok {
					// do-action, needs input from the table, but not found
					if !hasValue(v) {
						continue
					}
					nextArgs = mergeArgs(nextArgs, v)
				}
				nextExtra = mergeMap(nextExtra, p.findExrea(newExtra))
			} else { // search
				if !hasValue(nextArgs) && !hasValue(nextExtra) {
					continue
				}
			}
			newLists, err := m.RunContext(ctx, db, p.AtomName, p.ActionName, &RunOption{Args: nextArgs, Extra: nextExtra, GlobalArgs: globalArgs, GlobalExtra: globalExtra})
			if err != nil {
				return nil, err
			}
			if hasValue(newLists) {
				isRecursive := tableObj.IsRecursive()
				if isRecursive {
					// one-to-many recursive found
					short := p.shortenRecursive(newLists)
					item.(map[string]any)[p.Subname()] = short
				} else if tableObj.RecursiveColumn() != "" && pAtom.Table.IsRecursive() {
					switch p.Dimension {
					case CONNECTMap, CONNECTOne:
						item.(map[string]any)[p.Subname()] = newLists[0]
					default:
						item.(map[string]any)[p.Subname()] = newLists
					}
				} else if pAtom.Table.IsRecursive() {
					//short := ShortenX(p.Marker, newLists)
					short := p.shorten(newLists)
					item.(map[string]any)[p.Subname()] = short
				} else if tableObj.TableName == p.AtomName && p.Dimension == CONNECTOne {
					// simple loop table but not marked as isRecursive
					item.(map[string]any)[p.Subname()] = newLists[0]
				} else {
					short := p.shorten(newLists)
					item.(map[string]any)[p.Subname()] = short
				}
			}
		}
	}

	return data, nil
}
