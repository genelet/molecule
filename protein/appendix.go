package protein

import (
	"database/sql"
)

type Appendix struct {
	Addeds   []*Morph `json:"addeds,omitempty" hcl:"added,block"`
	Extras   []*Morph `json:"extras,omitempty" hcl:"extra,block"`
	DBAddeds []*Morph `json:"dbaddeds,omitempty" hcl:"db_added,block"`
	DBExtras []*Morph `json:"dbextras,omitempty" hcl:"db_extra,block"`
	Afters   []*Morph `json:"afters,omitempty" hcl:"after,block"`
}

type tmpAppendix struct {
	AppendObject *Appendix `json:"appendix,omitempty" hcl:"appendix,block"`
}

func (self *Appendix) RunBefore(db *sql.DB, args, extra map[string]interface{}) error {
	var err error
	if err = RunLoop(nil, self.Addeds, args, args); err == nil {
		if err = RunLoop(nil, self.Extras, args, extra); err == nil {
			if err = RunLoop(db, self.DBAddeds, args, args); err == nil {
				err = RunLoop(db, self.DBExtras, args, extra)
			}
		}
	}
	return err
}

func (self *Appendix) RunAfter(db *sql.DB, data *[]interface{}) error {
	for i, iitem := range *data {
		item := iitem.(map[string]interface{})
		err := RunLoop(nil, self.Afters, item, item)
		if err != nil {
			return err
		}
		(*data)[i] = item
	}

	return nil
}
