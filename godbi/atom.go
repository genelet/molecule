package godbi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/genelet/determined/dethcl"
	"github.com/genelet/determined/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Atom is a table with multiple actions
type Atom struct {
	AtomName string `json:"atomName,omitempty" hcl:"atomName,label"`
	Table
	Actions []Capability `json:"actions,omitempty" hcl:"actions,block"`
	customs map[string]any
}

// UnmarshalJSON is a JSON unmarshaler
func (self *Atom) UnmarshalJSON(bs []byte) error {
	type m struct {
		AtomName string `json:"atomName,omitempty"`
		Table
		Actions []map[string]any `json:"actions,omitempty"`
	}
	tmp := &m{}
	if err := json.Unmarshal(bs, tmp); err != nil {
		return err
	}

	trans := getEmptyCapacities()
	for _, action := range tmp.Actions {
		v, ok := action["actionName"]
		if !ok {
			continue
		}
		name := v.(string)
		jsonString, err := json.Marshal(action)
		if err != nil {
			return err
		}
		for i, item := range trans {
			if name == item.GetActionName() {
				if err = json.Unmarshal(jsonString, item); err != nil {
					return err
				}
				switch name {
				case "insert", "update", "insupd", "delete", "delecs":
					item.SetIsDo(true)
				default:
				}
				trans[i] = item
				break
			}
		}
	}

	self.AtomName = tmp.AtomName
	self.Table = tmp.Table
	self.Actions = trans
	return nil
}

// UnmarshalHCL is a HCL unmarshaler
func (self *Atom) UnmarshalHCL(bs []byte, labels ...string) error {
	file, diags := hclsyntax.ParseConfig(bs, rname(), hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return diags
	}
	spec, ref, err := specFromAtomBody(file.Body.(*hclsyntax.Body), self.customs)
	if err != nil {
		return err
	}

	err = dethcl.UnmarshalSpec(bs, self, spec, ref, labels...)
	if err != nil {
		return err
	}
	self.updateDefaultActions()
	return nil
}

func specFromAtomBody(body *hclsyntax.Body, customs map[string]any) (*utils.Struct, map[string]any, error) {
	ref := map[string]any{"Connection": new(Connection)}
	accepted := make(map[string]bool)
	for k, v := range customs {
		ref[k] = v
		accepted[k] = true
	}
	for _, v := range getEmptyCapacities() {
		ref[v.GetActionName()] = v
		accepted[v.GetActionName()] = true
	}

	var actions [][2]any
	for _, block := range body.Blocks {
		if block.Type != "actions" {
			continue
		}
		if len(block.Labels) == 0 {
			return nil, nil, fmt.Errorf("HCL finds no actionName from actions")
		}
		key := block.Labels[0]
		if _, ok := accepted[key]; !ok {
			continue
		}

		var nextpages, prepares []string
		for _, b := range block.Body.Blocks {
			if b.Type == "nextpages" {
				nextpages = append(nextpages, "Connection")
			}
			if b.Type == "prepares" {
				prepares = append(prepares, "Connection")
			}
		}
		second := make(map[string]any)
		if nextpages != nil {
			second["Nextpages"] = nextpages
		}
		if prepares != nil {
			second["Prepares"] = prepares
		}
		actions = append(actions, [2]any{key, second})
	}

	s, err := utils.NewStruct("Atom", map[string]any{"Actions": actions})
	return s, ref, err
}

func rname() string {
	return fmt.Sprintf("%d.hcl", rand.Int())
}

func getEmptyCapacities() []Capability {
	return []Capability{
		&Insert{Action: Action{ActionName: "insert"}},
		&Update{Action: Action{ActionName: "update"}},
		&Insupd{Action: Action{ActionName: "insupd"}},
		&Delete{Action: Action{ActionName: "delete"}},
		&Delecs{Action: Action{ActionName: "delecs"}},
		&Topics{Action: Action{ActionName: "topics"}},
		&Edit{Action: Action{ActionName: "edit"}},
		&Stmt{Action: Action{ActionName: "stmt"}},
	}
}

func (self *Atom) updateDefaultActions() {
	for _, v := range getEmptyCapacities() {
		found := false
		for _, existing := range self.Actions {
			if v.GetActionName() == existing.GetActionName() {
				found = true
				break
			}
		}
		if !found {
			self.Actions = append(self.Actions, v)
		}
	}
	for i, v := range self.Actions {
		switch v.GetActionName() {
		case "insert", "update", "insupd", "delete", "delecs":
			self.Actions[i].SetIsDo(true)
		default:
		}
	}
}

// MergeCustomActions merges custom actions
func (self *Atom) MergeCustomActions(custom ...Capability) {
	if custom == nil {
		return
	}

	names := make(map[string]int)
	for i, v := range self.Actions {
		names[v.GetActionName()] = i
	}

	clone := func(old any) any {
		obj := reflect.New(reflect.TypeOf(old).Elem())
		oldVal := reflect.ValueOf(old).Elem()
		newVal := obj.Elem()
		for i := 0; i < oldVal.NumField(); i++ {
			newValField := newVal.Field(i)
			if newValField.CanSet() {
				newValField.Set(oldVal.Field(i))
			}
		}
		return obj.Interface()
	}

	for _, v := range custom {
		actionName := v.GetActionName()
		if i, ok := names[actionName]; ok {
			self.Actions[i] = v
		} else {
			if self.customs == nil {
				self.customs = make(map[string]any)
			}
			self.customs[actionName] = clone(v)
			self.Actions = append(self.Actions, v)
		}
	}
}

// GetAction gets a specific action by name
func (self *Atom) GetAction(actionName string) Capability {
	for _, item := range self.Actions {
		if item.GetActionName() == actionName {
			return item
		}
	}

	return nil
}

// RunAtom runs an action by name
func (self *Atom) RunAtom(db *sql.DB, action string, ARGS any, extra ...map[string]any) ([]any, error) {
	return self.RunAtomContext(context.Background(), db, action, ARGS, extra...)
}

// RunAtomContext runs an action with context by name
func (self *Atom) RunAtomContext(ctx context.Context, db *sql.DB, action string, ARGS any, extra ...map[string]any) ([]any, error) {
	obj := self.GetAction(action)
	if obj == nil {
		return nil, errorActionNil(action)
	}
	if ARGS == nil {
		return obj.RunActionContext(ctx, db, &self.Table, nil, extra...)
	}

	switch t := ARGS.(type) {
	case map[string]any:
		return obj.RunActionContext(ctx, db, &self.Table, t, extra...)
	case []map[string]any:
		var data []any
		for _, item := range t {
			lists, err := obj.RunActionContext(ctx, db, &self.Table, item, extra...)
			if err != nil {
				return nil, err
			}
			data = append(data, lists...)
		}
		return data, nil
	case []any:
		var data []any
		for _, item := range t {
			if args, ok := item.(map[string]any); ok {
				lists, err := obj.RunActionContext(ctx, db, &self.Table, args, extra...)
				if err != nil {
					return nil, err
				}
				data = append(data, lists...)
			}
		}
		return data, nil
	default:
	}

	return nil, errorInputDataType(ARGS)
}
