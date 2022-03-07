package godbi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Stopper interface {
	Sign(tableObj *Table, item interface{}) bool
}

// Molecule describes all atoms and actions in a database schema
//
type Molecule struct {
	Atoms []Navigate `json:"atoms" hcl:"atoms"`
	DatabaseName string `json:"databaseName" hcl:"databaseName"`
	DBDriver DBType `json:"dbDriver" hcl:"dbDriver"`
	argsMap map[string]interface{}
	extraMap map[string]interface{}
	Stopper
}

func NewMoleculeJsonFile(fn string, cmap ...map[string][]Capability) (*Molecule, error) {
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return NewMoleculeJson(json.RawMessage(dat), cmap...)
}

type g struct {
	Atoms []json.RawMessage `json:"atoms" hcl:"atoms"`
	DatabaseName string     `json:"databaseName" hcl:"databaseName"`
	DBDriver DBType         `json:"dbDriver" hcl:"dbDriver"`
}

func NewMoleculeJson(dat json.RawMessage, cmap ...map[string][]Capability) (*Molecule, error) {
	tmps := new(g)
	if err := json.Unmarshal(dat, &tmps); err != nil {
		return nil, err
	}

	var atoms []Navigate
	for _, tmp := range tmps.Atoms {
		atom, err := NewAtomJson(tmp)
		if err != nil { return nil, err }
		if cmap != nil && cmap[0] != nil {
			cs := cmap[0][atom.Table.TableName]
			atom, err = NewAtomJson(tmp, cs...)
			if err != nil { return nil, err }
		}
		atoms = append(atoms, atom)
	}

	return &Molecule{Atoms:atoms, DatabaseName: tmps.DatabaseName, DBDriver:tmps.DBDriver}, nil
}

func (self *Molecule) String() string {
	bs, _ := json.MarshalIndent(self, "", "  ")
	return fmt.Sprintf("%s", bs)
}

func (self *Molecule) Initialize(args map[string]interface{}, extra map[string]interface{}) {
	self.argsMap = args
    self.extraMap = extra
}

func (self *Molecule) GetAtom(atom string) Navigate {
	if self.Atoms != nil {
		for _, item := range self.Atoms {
			tableObj := item.GetTable()
			if tableObj.GetTableName() == atom {
				tableObj.SetDBDriver(self.DBDriver)
				return item
			}
		}
	}
	return nil
}

// RunContext runs action by atom and action string names.
// It returns the searched data and optional error code.
//
// 'atom' is the atom name, and 'action' the action name.
// The first extra is the input data, shared by all sub actions.
// The rest are specific data for each action starting with the current one.
//
func (self *Molecule) RunContext(ctx context.Context, db *sql.DB, atom, action string, rest ...interface{}) ([]interface{}, error) {
	return self.generalContext(false, ctx, db, atom, action, rest...)
}

func (self *Molecule) runRecurseContext(ctx context.Context, db *sql.DB, atom, action string, rest ...interface{}) ([]interface{}, error) {
	return self.generalContext(true, ctx, db, atom, action, rest...)
}

func (self *Molecule) generalContext(topRecursive bool, ctx context.Context, db *sql.DB, atom, action string, rest ...interface{}) ([]interface{}, error) {
	var args interface{}
	var extra map[string]interface{}
	if rest != nil {
		if hasValue(rest[0]) {
			args = rest[0]
		}
		if len(rest) == 2 && hasValue(rest[1]) {
			switch t := rest[1].(type) {
			case map[string]interface{}: extra = t
			default:
				return nil, errorExtraDataType(rest[1])
			}
		}
	}

	if hasValue(self.argsMap[atom]) {
		argsMap := self.argsMap[atom].(map[string]interface{})
		args = MergeArgs(args, argsMap[action])
	}

	if hasValue(self.extraMap[atom]) {
		extraAction := self.extraMap[atom].(map[string]interface{})
		if hasValue(extraAction[action]) {
			extra = MergeMap(extra, extraAction[action].(map[string]interface{}))
		}
	}

	switch t := args.(type) {
	case map[string]interface{}:
		return self.hashContext(topRecursive, ctx, db, atom, action, t, extra)
	case []map[string]interface{}:
		var final []interface{}
		for _, arg := range t {
			lists, err := self.hashContext(topRecursive, ctx, db, atom, action, arg, extra)
			if err != nil { return nil, err }
			final = append(final, lists...)
		}
		return final, nil
	case []interface{}:
		var final []interface{}
		for _, arg := range t {
			if v, ok := arg.(map[string]interface{}); ok {
				lists, err := self.hashContext(topRecursive, ctx, db, atom, action, v, extra)
				if err != nil { return nil, err }
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

func (self *Molecule) hashContext(topRecursive bool, ctx context.Context, db *sql.DB, atom, action string, args, extra map[string]interface{}) ([]interface{}, error) {
	atomObj := self.GetAtom(atom)
	if atomObj == nil {
		return nil, errorAtomNotFound(atom)
	}
	tableObj := atomObj.GetTable()
	actionObj := atomObj.GetAction(action)
	if actionObj == nil {
		return nil, errorActionNotFound(action, atom)
	}
	if actionObj.GetIsDo() && !hasValue(args) {
		return nil, nil
	}

	prepares := actionObj.GetPrepares()
	nextpages := actionObj.GetNextpages()

	newArgs := CloneArgs(args)
	newExtra := CloneMap(extra)

	// prepares receives filtered args and extra from current args
	for _, p := range prepares {
		pAtom       := self.GetAtom(p.TableName)
		pTable      := pAtom.GetTable()
		v, _        := p.FindArgs(args)
		preArgs     := MergeArgs(p.NextArgs(args), v)
		preExtra    := MergeMap(p.NextExtra(args), p.FindExtra(extra))
		isDo        := pAtom.GetAction(p.ActionName).GetIsDo()
		isRecursive := pTable.IsRecursive()

		var lists []interface{}
		var err error

		if topRecursive {
			if !hasValue(preArgs) {
				return []interface{}{args}, nil
			}
			pk := pTable.Pks[0]
			if isRecursive {
				switch t := preArgs.(type) {
				case map[string]interface{}: delete(t, pk)
				case []interface{}:
                for _, s := range t { delete(s.(map[string]interface{}), pk) }
				case []map[string]interface{}:
                for _, s := range t { delete(s, pk) }
				}
			}
			lists, err = self.runRecurseContext(ctx, db, p.TableName, p.ActionName, preArgs, preExtra)
		} else if isDo && isRecursive {
		// this triggers the original topRecursive and is always a DO action
			lists, err = self.runRecurseContext(ctx, db, p.TableName, p.ActionName, preArgs, preExtra)
		} else {
			lists, err = self.RunContext(ctx, db, p.TableName, p.ActionName, preArgs, preExtra)
		}
		if err != nil { return nil, err }

		// only two types of prepares
		// 1) one pre, with multiple outputs (when p.argsMap is multiple)
		if hasValue(lists) && len(lists) > 1 {
			var tmp []map[string]interface{}
			newExtra = CloneMap(extra)
			for _, item := range lists {
				result := MergeArgs(args, p.NextArgs(item)).(map[string]interface{})
				tmp = append(tmp, result)
				newExtra = MergeMap(newExtra, p.NextExtra(item))
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
			middle := p.NextArgs(lists[0])
			if topRecursive && isRecursive {
				pk := pTable.Pks[0]
				switch t := middle.(type) {
				case map[string]interface{}:
					delete(t, pk)
				case []interface{}:
					for _, s := range t {
						delete(s.(map[string]interface{}), pk)
					}
				case []map[string]interface{}:
					for _, s := range t { delete(s, pk) }
				}
			}
			newArgs = MergeArgs(newArgs, middle.(map[string]interface{}), true)
			newExtra = MergeMap(newExtra, p.NextExtra(lists[0]))
		}
	}

	if !topRecursive && hasValue(newArgs) && actionObj.GetIsDo() {
		//newArgs = tableObj.refreshArgs(newArgs, true)
		newArgs = tableObj.refreshArgs(newArgs)
	}

	data, err := atomObj.RunAtomContext(ctx, db, action, newArgs, newExtra)
	if err != nil { return nil, err }

	if topRecursive && action != "delecs" {
		if tableObj.IsRecursive() {
			var p      *Connection
			var pAtom   Navigate
			var rColumn string
			for _, p = range nextpages {
				pAtom = self.GetAtom(p.TableName)
				rColumn = pAtom.GetTable().RecursiveColumn()
				if rColumn != "" { break }
			}
			if rColumn == "" { return data, nil }
			argsData, _ := p.FindArgs(newArgs)
			if !hasValue(argsData) {
				return data, nil
			}
			for _, item := range data {
				nextArgs :=  p.NextArgs(item)
				nextArgs = MergeArgs(argsData, nextArgs)
				_, err = self.runRecurseContext(ctx, db, p.TableName, p.ActionName, nextArgs, nil)
				if err != nil { return nil, err }
			}
		}
		return data, err
	}

	for _, p := range nextpages {
		for _, item := range data {
			if self.Stopper != nil && self.Stopper.Sign(tableObj, item) { continue }
			pAtom := self.GetAtom(p.TableName)
			pAction := pAtom.GetAction(p.ActionName)
			nextArgs := p.NextArgs(item)
			nextExtra := p.NextExtra(item)
			if pAction.GetIsDo() {
				if v, ok := p.FindArgs(newArgs); ok {
					// do-action, needs input from the table, but not found
					if !hasValue(v) {
						continue
					}
					nextArgs  = MergeArgs(nextArgs, v)
				}
				nextExtra = MergeMap(nextExtra, p.FindExtra(newExtra))
			} else { // search
				if !hasValue(nextArgs) && !hasValue(nextExtra) {
					continue
				}
			}
			newLists, err := self.RunContext(ctx, db, p.TableName, p.ActionName, nextArgs, nextExtra)
			if err != nil { return nil, err }
			if hasValue(newLists) {
				isRecursive := tableObj.IsRecursive()
				if isRecursive {
					// one-to-many recursive found
					short := p.ShortenRecursive(newLists)
					item.(map[string]interface{})[p.Subname()] = short
				} else if tableObj.RecursiveColumn() != "" && pAtom.GetTable().IsRecursive() {
					switch p.Dimension {
					case CONNECTMap, CONNECTOne:
						item.(map[string]interface{})[p.Subname()] = newLists[0]
					default:
						item.(map[string]interface{})[p.Subname()] = newLists
					}
				} else if pAtom.GetTable().IsRecursive() {
					//short := ShortenX(p.Marker, newLists)
					short := p.Shorten(newLists)
					item.(map[string]interface{})[p.Subname()] = short
				} else if tableObj.TableName == p.TableName && p.Dimension == CONNECTOne {
					// simple loop table but not marked as isRecursive
					item.(map[string]interface{})[p.Subname()] = newLists[0]
				} else {
					short := p.Shorten(newLists)
					item.(map[string]interface{})[p.Subname()] = short
				}
			}
		}
	}

	return data, nil
}
