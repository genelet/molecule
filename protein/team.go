package protein

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/genelet/molecule/godbi"
	"github.com/genelet/team/micro"
)

type Team struct {
	Appendix map[string]interface{} `json:"appendix,omitempty"`

	UserTable  string          `json:"userid_table,omitempty"`
	UserIDName string          `json:"userid_name,omitempty"`
	IsAdmin    bool            `json:"is_admin,omitempty"`
	IsPublic   bool            `json:"is_public,omitempty"`
	Molecule   *godbi.Molecule `json:"colorfuls,omitempty"`
}

func (self *Team) GetColorful(tableName string) *Colorful {
	for _, colorful := range self.Colorfuls {
		if tableName == colorful.Atom.GetTable().TableName {
			return colorful
		}
	}

	return nil
}

func (self *Team) String() string {
	bs, _ := json.MarshalIndent(self, "", "  ")
	return fmt.Sprintf("%s", bs)
}

func (self *Team) Dashboard() ([]string, []string) {
	var generics, landings []string
	for _, colorful := range self.Colorfuls {
		tableName := colorful.Atom.Table.TableName
		if !colorful.IsProtect() {
			generics = append(generics, tableName)
		} else if colorful.IsLanding {
			landings = append(landings, tableName)
		}
	}
	return generics, landings
}

func (self *Team) RunContext(ctx context.Context, db *sql.DB, dbname string, driver godbi.DBType, token, atom, action string, args, extra map[string]interface{}) ([]interface{}, error) {
	colorful := self.GetColorful(atom)
	if colorful == nil {
		return nil, fmt.Errorf("atom %s not found", atom)
	}

	if token != "" {
		err := colorful.checkFK(self.UserIDName, token, action, args, extra)
		if err != nil {
			return nil, err
		}
	}

	appendix, err := colorful.GetAppendix(action)
	if err != nil {
		return nil, err
	}

	if appendix != nil {
		err = appendix.RunBefore(db, args, extra)
		if err != nil {
			return nil, fmt.Errorf("before error in model %s and action %s: %v", atom, action, err)
		}
	}

	molecule := self.toMolecule(dbname, driver, &OneStepStopper{atom, token})
	data, err := molecule.RunContext(ctx, db, atom, action, args, extra)
	if err != nil {
		return nil, fmt.Errorf("graph error in model %s and action %s: %v", atom, action, err)
	}

	if data != nil && len(data) > 0 && appendix != nil {
		if err = appendix.RunAfter(db, &data); err != nil {
			return nil, fmt.Errorf("after error in model %s and action %s: %v", atom, action, err)
		}
	}

	return data, nil
}

func checkFK(colorful *Colorful, pk, token, action string, args, extra map[string]interface{}) error {
	if !colorful.IsProtect() {
		return nil
	}

	isDo := colorful.Atom.GetAction(action).GetIsDo()

	if (colorful.IsUser || colorful.IsLanding) && isDo == false {
		extra[pk] = args[pk]
		return nil
	}

	fk := colorful.ProtectedFk()
	if fk == nil {
		return fmt.Errorf("protected key not defined")
	}

	fkColumn := args[fk.FkColumn]
	fkMd5 := args[fk.FkColumn+"_sign"]
	if fkColumn != nil && fkMd5 != nil {
		if digest(token, fk.FkTable, fk.FkColumn, fmt.Sprintf("%v", fkColumn)) == fkMd5.(string) {
			if isDo {
				args[fk.Column] = fkColumn
			} else {
				extra[fk.Column] = fkColumn
			}
			return nil
		}
		return fmt.Errorf("signature failed")
	}

	return fmt.Errorf("signature not found")
}

func digest(a1, a2, a3, a4 string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(a2+a1+a3+a4)))
}
