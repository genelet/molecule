package godbi

import (
	"context"
	"database/sql"
)

// Update struct for row update by primary key
type Update struct {
	Action
	Empties []string `json:"empties,omitempty" hcl:"empties,optional"`
}

var _ Capability = (*Update)(nil)

func (self *Update) RunAction(db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Update) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	if err := t.checkNull(ARGS); err != nil {
		return nil, err
	}

	ids := t.getIDVal(ARGS)
	if !hasValue(ids) {
		return nil, errorMissingPk(t.TableName)
	}

	fieldValues, allAuto := t.getFv(ARGS, self.getAllowed())
	if allAuto {
		return fromFv(fieldValues), nil
	}
	if !hasValue(fieldValues) {
		return nil, errorEmptyInput(t.TableName)
	} else if len(fieldValues) == 1 && t.Pks != nil {
		for _, pk := range t.Pks {
			if _, ok := fieldValues[pk]; ok {
				return fromFv(fieldValues), nil
			}
		}
		return fromFv(fieldValues), nil
	}

	err := t.updateHashNullsContext(ctx, db, fieldValues, ids, self.Empties, extra...)
	return fromFv(fieldValues), err
}
