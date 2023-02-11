package protein

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/genelet/molecule/godbi"
)

type Team struct {
	UserIDName string        `json:"userid_name,omitempty"`
	IsAdmin    bool          `json:"is_admin,omitempty"`
	IsPublic   bool          `json:"is_public,omitempty"`
	Colorfuls  []*Colorful   `json:"colorfuls,omitempty"`
}

func (self *Team) String() string {
	bs, _ := json.MarshalIndent(self, "", "  ")
	return fmt.Sprintf("%s", bs)
}

func (self *Team) toMolecule(driver godbi.DBType, stopper godbi.Stopper) *godbi.Molecule {
	var atoms []godbi.Navigate
	for _, colorful := range self.Colorfuls {
		atoms = append(atoms, colorful.Atom)
	}
	return &godbi.Molecule{
		Atoms:        atoms,
		DBDriver:     driver,
		Stopper:      stopper}
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

func (self *Team) RunContext(ctx context.Context, db *sql.DB, driver godbi.DBType, token, atom, action string, args, extra map[string]interface{}) ([]interface{}, error) {
	var colorful *Colorful
	for _, c := range self.Colorfuls {
		if atom == c.Atom.GetTable().TableName {
			colorful = c
			break
		}
	}
	if colorful == nil {
		return nil, fmt.Errorf("colorful not found for %s", atom)
	}

	if !self.IsAdmin && !self.IsPublic {
		if token == "" {
			return nil, fmt.Errorf("token mising in protected team %s", atom)
		}
		err := checkFK(colorful, self.UserIDName, token, action, args, extra)
		if err != nil {
			return nil, err
		}
	}

	var appendix *Appendix
	if colorful.Appendices != nil {
		appendix = colorful.Appendices[action]
		err := appendix.RunBefore(db, args, extra)
		if err != nil {
			return nil, fmt.Errorf("before error in model %s and action %s: %v", atom, action, err)
		}
	}

	molecule := self.toMolecule(driver, &OneStepStopper{atom, token})
	data, err := molecule.RunContext(ctx, db, atom, action, args, extra)
	if err == nil && data != nil && len(data) > 0 && appendix != nil {
		err = appendix.RunAfter(db, &data)
	}

	return data, err
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

	fk := colorful.protectedFk()
	if fk == nil {
		return fmt.Errorf("protected key not defined")
	}

	fkColumn := args[fk.FkColumn]
	fkMd5 := args[fk.FkColumn+SignPOSTFIX]
	if fkColumn != nil && fkMd5 != nil {
		if digest(token, fk.FkTable, fk.FkColumn, fkColumn) == fkMd5.(string) {
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

func digest(a1, a2, a3 string, a4 interface{}) string {
	str := fmt.Sprintf("%v", a4)
	return fmt.Sprintf("%x", sha1.Sum([]byte(a2+a1+a3+str)))
}
