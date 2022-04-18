package godbi

import (
	"context"
	"database/sql"
)

// Insert struct for table insert
type Insert struct {
	Action
}

// Run inserts a row using data passed in ARGS.
//
func (self *Insert) RunAction(db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

// InsertContext inserts a row using data passed in ARGS.
//
func (self *Insert) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	if self.IsDo {
		if err := t.checkNull(ARGS); err != nil {
			return nil, err
		}
	}

	fieldValues, allAuto := t.getFv(ARGS)
	if !allAuto && !hasValue(fieldValues) {
		return nil, errorEmptyInput(t.TableName)
	}

	autoID, err := t.insertHashContext(ctx, db, fieldValues)
	if err != nil {
		return nil, err
	}

	if t.IdAuto != "" {
		fieldValues[t.IdAuto] = autoID
	}

	return fromFv(fieldValues), nil
}
