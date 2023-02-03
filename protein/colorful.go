package protein

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/genelet/molecule/godbi"
	"github.com/genelet/team/micro"
)

const (
	SIGN = "_sign"
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
	if self.Protected != "" {
		return true
	}
	return false
}

type Colorful struct {
	Color `json:"color,omitempty" hcl:"color,block"`
	Atom  *godbi.Atom `json:"atom,omitempty" hcl:"atom,block"`
}

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

	if !self.IsUser && self.IsLanding && self.Color.Protected == "" {
		return self.Atom.RunAtomContext(ctx, db, action, ARGS, extra...)
	}

	if self.whoPublic {
		return nil, fmt.Errorf("protected area")
	}

	if !self.whoAdmin {
		var args0, extra0 map[string]interface{}
		if extra != nil {
			extra0 = extra[0]
		} else {
			extra0 = make(map[string]interface{})
		}
		if ARGS != nil {
			args0 = ARGS.(map[string]interface{})
		}
        err := self.checkFK(action, args0, extra0)
        if err != nil {
            return nil, err
        }
    }

	return self.Atom.RunAtomContext(ctx, db, action, ARGS, extra...)
}

func (self *Colorful) checkFK(pk, token, action string, args, extra map[string]interface{}) error {
    isDo := self.GetAction(action).GetIsDo()

    if (self.IsUser || self.IsLanding) && isDo == false {
        extra[pk] = args[pk]
        return nil
    }

	table := self.GetTable() {
	var fkObj *godib.Fk
	for _, fk := range table.Fks {
		if fk.Column == self.Protected {
			fkObj = fk
			break // the same column could be protected by multiple fks?
		}
	}
    if fk == nil {
        return fmt.Errorf("protected column in foreign table not found")
    }

    fkColumn := args[fkObj.FkColumn]
    fkMd5 := args[fkObj.FkColumn+SIGN]
    if fkColumn != nil && fkMd5 != nil {
        if digest(token, fkObj.FkTable, fkObj.FkColumn, fkColumn) == fkMd5.(string) {
            if isDo {
                args[fkObj.Column] = fkColumn
            } else {
                extra[fkObj.Column] = fkColumn
            }
            return nil
        }
        return fmt.Errorf("signature failed")
    }

    return fmt.Errorf("signature not found")
}

func digest(a1, a2, a3 string, v interface{}) string {
	a4 := fmt.Sprintf("%v", v)
    return fmt.Sprintf("%x", sha1.Sum([]byte(a2+a1+a3+a4)))
}
