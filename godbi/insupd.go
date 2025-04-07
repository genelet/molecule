package godbi

import (
	"context"
	"database/sql"
)

// Insupd struct for table update, if not existing according to the unique key, do insert
type Insupd struct {
	Action
}

var _ Capability = (*Insupd)(nil)

func (self *Insupd) RunAction(db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Insupd) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	if err := t.checkNull(ARGS); err != nil {
		return nil, err
	}

	fieldValues, allAuto := t.getFv(ARGS, self.getAllowed())
	if !allAuto && !hasValue(fieldValues) {
		return nil, errorEmptyInput(t.TableName)
	}

	changed, err := t.insupdTableContext(ctx, db, fieldValues)
	if err != nil {
		return nil, err
	}

	if t.IDAuto != "" {
		fieldValues[t.IDAuto] = changed
	}

	return fromFv(fieldValues), nil
}
