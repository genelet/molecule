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
	logger   Slogger
	argsMap  map[string]any
	extraMap map[string]any
}

// SetLogger sets the logger
func (self *Molecule) SetLogger(logger Slogger) {
	self.logger = logger
	for _, atom := range self.Atoms {
		atom.Table.SetLogger(logger)
	}
}

// GetLogger gets the logger
func (self *Molecule) GetLogger() Slogger {
	return self.logger
}

// Initialize initializes molecule with args and extra
// args is the input data, shared by all sub actions.
// extra is specific data list for sub actions, starting with the current one.
func (self *Molecule) Initialize(args map[string]any, extra map[string]any) {
	self.argsMap = args
	self.extraMap = extra
}

// GetAtom returns the atom by atom name
func (self *Molecule) GetAtom(atomName string) *Atom {
	if self.Atoms != nil {
		for _, atom := range self.Atoms {
			if atom.AtomName == atomName {
				atom.SetDBDriver(self.DBDriver)
				atom.Table.logger = self.logger
				return atom
			}
		}
	}
	return nil
}

// RunContext runs action by atom and action string names.
// It returns the searched data and optional error code.
// atom is the atom name, and action the action name. rest are:
//   - the input data, shared by all sub actions.
//   - specific data list for sub actions, starting with the current one.
func (self *Molecule) RunContext(ctx context.Context, db *sql.DB, atom, action string, rest ...any) ([]any, error) {
	return self.generalContext(false, ctx, db, atom, action, rest...)
}

func (self *Molecule) runRecurseContext(ctx context.Context, db *sql.DB, atom, action string, rest ...any) ([]any, error) {
	return self.generalContext(true, ctx, db, atom, action, rest...)
}

func (self *Molecule) generalContext(topRecursive bool, ctx context.Context, db *sql.DB, atom, action string, rest ...any) ([]any, error) {
	var args any
	var extra map[string]any
	if rest != nil {
		if hasValue(rest[0]) {
			args = rest[0]
		}
		if len(rest) == 2 && hasValue(rest[1]) {
			switch t := rest[1].(type) {
			case map[string]any:
				extra = t
			default:
				return nil, errorExtraDataType(rest[1])
			}
		}
	}

	if hasValue(self.argsMap[atom]) {
		argsMap := self.argsMap[atom].(map[string]any)
		args = mergeArgs(args, argsMap[action])
	}

	if hasValue(self.extraMap[atom]) {
		extraAction := self.extraMap[atom].(map[string]any)
		if hasValue(extraAction[action]) {
			extra = mergeMap(extra, extraAction[action].(map[string]any))
		}
	}

	switch t := args.(type) {
	case map[string]any:
		return self.hashContext(topRecursive, ctx, db, atom, action, t, extra)
	case []map[string]any:
		var final []any
		for _, arg := range t {
			lists, err := self.hashContext(topRecursive, ctx, db, atom, action, arg, extra)
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
				lists, err := self.hashContext(topRecursive, ctx, db, atom, action, v, extra)
				if err != nil {
					return nil, err
				}
				final = append(final, lists...)
			}
		}
		return final, nil
	default:
	}

	return self.hashContext(topRecursive, ctx, db, atom, action, nil, extra)
}

// RunContext runs action by atom and action string names.
// It returns the searched data and optional error code.
//
// 'atom' is the atom name, and 'action' the action name.
// The first extra is the input data, shared by all sub actions.
// The rest are specific data for each action starting with the current one.
//

func (self *Molecule) hashContext(topRecursive bool, ctx context.Context, db *sql.DB, atom, action string, args, extra map[string]any) ([]any, error) {
	atomObj := self.GetAtom(atom)
	if atomObj == nil {
		return nil, errorAtomNotFound(atom)
	}
	tableObj := atomObj.Table
	actionObj := atomObj.GetAction(action)
	if actionObj == nil {
		return nil, errorActionNotFound(action, atom)
	}
	if actionObj.GetIsDo(action) && !hasValue(args) {
		return nil, nil
	}

	prepares := actionObj.GetPrepares()
	nextpages := actionObj.GetNextpages()

	newArgs := cloneArgs(args)
	newExtra := cloneMap(extra)

	// prepares receives filtered args and extra from current args
	for _, p := range prepares {
		pAtom := self.GetAtom(p.AtomName)
		pTable := pAtom.Table
		v, _ := p.findArgs(args)
		preArgs := mergeArgs(p.nextArgs(args), v)
		preExtra := mergeMap(p.nextExtra(args), p.findExrea(extra))
		isDo := pAtom.GetAction(p.ActionName).GetIsDo(p.ActionName)
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
			lists, err = self.runRecurseContext(ctx, db, p.AtomName, p.ActionName, preArgs, preExtra)
		} else if isDo && isRecursive {
			// this triggers the original topRecursive and is always a DO action
			lists, err = self.runRecurseContext(ctx, db, p.AtomName, p.ActionName, preArgs, preExtra)
		} else if self.PreStopper == nil || !self.PreStopper.Stop(&tableObj, &pTable) {
			lists, err = self.RunContext(ctx, db, p.AtomName, p.ActionName, preArgs, preExtra)
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

	if !topRecursive && hasValue(newArgs) && actionObj.GetIsDo(action) {
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
		if self.Stopper != nil {
			self.Stopper.Sign(&tableObj, hash)
		}
	}

	if topRecursive && action != "delecs" {
		if tableObj.IsRecursive() {
			var p *Connection
			var pAtom *Atom
			var rColumn string
			for _, p = range nextpages {
				pAtom = self.GetAtom(p.AtomName)
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
				_, err = self.runRecurseContext(ctx, db, p.AtomName, p.ActionName, nextArgs, nil)
				if err != nil {
					return nil, err
				}
			}
		}
		return data, err
	}

	for _, p := range nextpages {
		pAtom := self.GetAtom(p.AtomName)
		pAction := pAtom.GetAction(p.ActionName)
		if self.Stopper != nil && self.Stopper.Stop(&tableObj, &(pAtom.Table)) {
			continue
		}

		for _, item := range data {
			if item == nil {
				continue
			}
			nextArgs := p.nextArgs(item)
			nextExtra := p.nextExtra(item)
			if pAction.GetIsDo(p.ActionName) {
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
			newLists, err := self.RunContext(ctx, db, p.AtomName, p.ActionName, nextArgs, nextExtra)
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
