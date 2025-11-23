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

// RunAction inserts a row using data passed in args.
func (i *Insert) RunAction(db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	return i.RunActionContext(context.Background(), db, t, args, extra...)
}

// RunActionContext inserts a row using data passed in args.
func (i *Insert) RunActionContext(ctx context.Context, db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	if err := t.checkNull(args); err != nil {
		return nil, err
	}

	fieldValues, allAuto := t.getFv(args, i.getAllowed())
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
