package protein

import (
	"fmt"

	"github.com/genelet/molecule/godbi"
)

type Color struct {
	IsUser    bool   `json:"is_user,omitempty" hcl:"is_user,optional"`
	IsLanding bool   `json:"is_landing,omitempty" hcl:"is_landing,optional"`
	Protected string `json:"protected,omitempty" hcl:"protected,optional"`
}

func NewColor(isUser, isLanding bool, protected string) *Color {
	return &Color{IsUser: isUser, Protected: protected, IsLanding: isLanding}
}

func (self *Color) IsProtect() bool {
	if self.IsUser || self.IsLanding || self.Protected != "" {
		return true
	}
	return false
}

type Colorful struct {
	Color                           `json:"color,omitempty" hcl:"color,block"`
	Atom       *godbi.Atom          `json:"atom,omitempty" hcl:"atom,block"`
	Appendices map[string]*Appendix `json:"appendices,omitempty" hcl:"appendices,block"`
}

/*
var _ godbi.Navigate = (*Colorful)(nil)

func (self *Colorful) GetTable() *godbi.Table {
    return self.Atom.GetTable()
}

func (self *Colorful) GetAction(action string) Capability {
    return self.Atom.GetAction(action)
}

func (self *Colorful) RunAtomContext(ctx context.Context, db *sql.DB, action string, ARGS interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.Atom.RunAtomContext(ctx, db, action, ARGS, extra...)
}
*/

func (self *Colorful) protectedFk() *godbi.Fk {
	for _, fk := range self.Atom.Table.Fks {
		if fk.Column == self.Protected {
			return fk
		}
	}
	return nil
}

func (self *Colorful) PublicPathToArgs(action string, paths []string, args map[string]interface{}) (string, error) {
    newAction := action
    n := len(paths)

    if n%2 == 1 {
        args[self.Atom.GetTable().Pks[0]] = paths[0]
        paths = paths[1:]
        n--
    } else if action == "edit" { // even with get is always list
        newAction = "topics"
    }

    for i := 0; i < n/2; i++ {
        args[paths[i*2]] = paths[i*2+1]
    }

    return newAction, nil
}

func (self *Colorful) ProtectPathToArgs(action string, paths []string, args map[string]interface{}) (string, error) {
    newAction := action
    n := len(paths)

    if self.IsUser {
        if action == "topics" { // absolutely no topics
            return "", fmt.Errorf("list action is disabled for user component")
        }
        if n%2 == 1 {
            return "", fmt.Errorf("wrong request for user component")
        }
        for i := 0; i < n/2; i++ {
            args[paths[i*2]] = paths[i*2+1]
        }
        return newAction, nil
    }

    table := self.Atom.GetTable()

    if n%2 == 1 {
        args[table.Pks[0]] = paths[0]
        paths = paths[1:]
        n--
    } else if action == "edit" { // even with get is always list
        newAction = "topics"
    }

   if !self.IsLanding && (action == "edit" || action == "topics") {
        if n < 2 {
            return "", fmt.Errorf("signature not found in the url path")
        } else if !self.IsProtect() {
            return "", fmt.Errorf("foreign key required but not defined")
        } else {
            args[self.Protected] = paths[0]
            args[self.Protected+SignPOSTFIX] = paths[1]
            paths = paths[2:]
            n--
            n--
        }
    }
    for i := 0; i < n/2; i++ {
        args[paths[i*2]] = paths[i*2+1]
    }

    return newAction, nil
}
