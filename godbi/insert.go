package godbi

import (
	"context"
	"database/sql"
)

// Insert struct for table insert
type Insert struct {
	Action
}

var _ Capability = (*Insert)(nil)

// RunAction inserts a row using data passed in ARGS.
func (self *Insert) RunAction(db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

// RunActionContext inserts a row using data passed in ARGS.
func (self *Insert) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	if err := t.checkNull(ARGS); err != nil {
		return nil, err
	}

	fieldValues, allAuto := t.getFv(ARGS, self.getAllowed())
	if !allAuto && !hasValue(fieldValues) {
		return nil, errorEmptyInput(t.TableName)
	}

	autoID, err := t.insertHashContext(ctx, db, fieldValues)
	if err != nil {
		return nil, err
	}

	if t.IDAuto != "" {
		fieldValues[t.IDAuto] = autoID
	}

	return fromFv(fieldValues), nil
}
