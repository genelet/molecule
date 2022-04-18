package godbi

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
)

// Atom is to implement Navigate interface
//
type Navigate interface {
	GetTable() *Table
	GetAction(string) Capability
	RunAtomContext(context.Context, *sql.DB, string, interface{}, ...map[string]interface{}) ([]interface{}, error)
}

// Atom is a table with multiple actions
type Atom struct {
	Table
	Actions []Capability `json:"actions,omitempty" hcl:"actions,optional"`
}

// NewAtomJsonFile parse a disk file to atom
func NewAtomJsonFile(fn string, custom ...Capability) (*Atom, error) {
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return NewAtomJson(json.RawMessage(dat), custom...)
}

type m struct {
	Table
	Actions []interface{} `json:"actions,omitempty" hcl:"actions,optional"`
}

// Atom is an unmarshaler
func (self *Atom) UnmarshalJSON(bs []byte) error {
	tmp := &m{}
	if err := json.Unmarshal(bs, tmp); err != nil {
		return err
	}
	actions, err := assertion(tmp.Actions)
	if err != nil {
		return err
	}
	self.Table = tmp.Table
	self.Actions = actions
	return nil
}

// NewAtomJson parses bytes to atom
func NewAtomJson(dat json.RawMessage, custom ...Capability) (*Atom, error) {
	tmp := &m{}
	if err := json.Unmarshal(dat, tmp); err != nil {
		return nil, err
	}
	actions, err := assertion(tmp.Actions, custom...)
	return &Atom{tmp.Table, actions}, err
}

func assertion(actions []interface{}, custom ...Capability) ([]Capability, error) {
	var trans []Capability

	for _, item := range actions {
		action := item.(map[string]interface{})
		v, ok := action["actionName"]
		if !ok { continue }
		name := v.(string)
		jsonString, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		var tran Capability
		found := false
		for _, item := range custom {
			if name==item.GetActionName() {
				tran = item
				found = true
				break
			}
		}
		if !found {
			switch name {
			case "insert":
				tran = &Insert{Action:Action{IsDo:true}}
			case "update":
				tran = &Update{Action:Action{IsDo:true}}
			case "insupd":
				tran = &Insupd{Action:Action{IsDo:true}}
			case "edit":
				tran = new(Edit)
			case "topics":
				tran = new(Topics)
			case "delete":
				tran = &Delete{Action:Action{IsDo:true}}
			case "delecs":
				tran = &Delecs{Action:Action{IsDo:true}}
			default:
				return nil, errorActionNotDefined(name)
			}
		}
		if err := json.Unmarshal(jsonString, tran); err != nil {
			return nil, err
		}
		trans = append(trans, tran)
	}
	return trans, nil
}

// GetTable gets the table from atom
func (self *Atom) GetTable() *Table {
	return &self.Table
}

// GetAction gets a specific action by name
func (self *Atom) GetAction(action string) Capability {
	if self.Actions != nil {
		for _, item := range self.Actions {
			if item.GetActionName() == action {
				return item
			}
		}
	}

	return nil
}

// RunAtom runs an action by name
func (self *Atom) RunAtom(db *sql.DB, action string, ARGS interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunAtomContext(context.Background(), db, action, ARGS, extra...)
}

// RunAtom runs an action with context by name
func (self *Atom) RunAtomContext(ctx context.Context, db *sql.DB, action string, ARGS interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
    obj := self.GetAction(action)
    if obj == nil {
        return nil, errorActionNil(action)
    }
	if ARGS == nil {
		return obj.RunActionContext(ctx, db, &self.Table, nil, extra...)
	}

	switch t := ARGS.(type) {
	case map[string]interface{}:
		return obj.RunActionContext(ctx, db, &self.Table, t, extra...)
	case []map[string]interface{}:
		var data []interface{}
		for _, item := range t {
			lists, err := obj.RunActionContext(ctx, db, &self.Table, item, extra...)
			if err != nil {
				return nil, err
			}
			data = append(data, lists...)
		}
		return data, nil
	case []interface{}:
		var data []interface{}
		for _, item := range t {
			if args, ok := item.(map[string]interface{}); ok {
				lists, err := obj.RunActionContext(ctx, db, &self.Table, args, extra...)
				if err != nil {
					return nil, err
				}
				data = append(data, lists...)
			}
		}
		return data, nil
	default:
		return nil, errorInputDataType(t)
	}

	return nil, nil
}
