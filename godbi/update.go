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

func (u *Update) RunAction(db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	return u.RunActionContext(context.Background(), db, t, args, extra...)
}

func (u *Update) RunActionContext(ctx context.Context, db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	if err := t.checkNull(args); err != nil {
		return nil, err
	}

	ids := t.getIDVal(args)
	if !hasValue(ids) {
		return nil, errorMissingPk(t.TableName)
	}

	fieldValues, allAuto := t.getFv(args, u.getAllowed())
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

	err := t.updateHashNullsContext(ctx, db, fieldValues, ids, u.Empties, extra...)
	return fromFv(fieldValues), err
}
