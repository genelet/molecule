package molecule

import (
	"context"
	"database/sql"
)

type Insupd struct {
	Action
}

func (self *Insupd) RunAction(db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Insupd) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	if self.IsDo {
		if err := t.checkNull(ARGS); err != nil {
			return nil, err
		}
	}

	fieldValues, allAuto := t.getFv(ARGS)
	if !allAuto && !hasValue(fieldValues) {
		return nil, errorEmptyInput(t.TableName)
	}

	changed, err := t.insupdTableContext(ctx, db, fieldValues)
	if err != nil {
		return nil, err
	}

	if t.IdAuto != "" {
		fieldValues[t.IdAuto] = changed
	}

	return fromFv(fieldValues), nil
}
