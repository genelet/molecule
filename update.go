package molecule

import (
	"context"
	"database/sql"
)

type Update struct {
	Action
	Empties []string `json:"empties,omitempty" hcl:"empties,optional"`
}

func (self *Update) RunAction(db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Update) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	if self.IsDo {
		if err := t.checkNull(ARGS); err != nil {
			return nil, err
		}
	}

	ids := t.getIdVal(ARGS)
	if !hasValue(ids) {
		return nil, errorMissingPk(t.TableName)
	}

	fieldValues, allAuto := t.getFv(ARGS)
	if allAuto {
		return fromFv(fieldValues), nil
	}
	if !hasValue(fieldValues) {
		return nil, errorEmptyInput(t.TableName)
	} else if len(fieldValues) == 1 && fieldValues[t.Pks[0]] != nil {
		return fromFv(fieldValues), nil
	}

	err := t.updateHashNullsContext(ctx, db, fieldValues, ids, self.Empties, extra...)
	return fromFv(fieldValues), err
}
