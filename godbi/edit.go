package godbi

import (
	"context"
	"database/sql"
)

// Edit struct for search one specific row by primary key
type Edit struct {
	Action
	FIELDS string `json:"fields,omitempty" hcl:"fields,optional"`
}

var _ Capability = (*Edit)(nil)

func (e *Edit) setDefaultElementNames() []string {
	if e.FIELDS == "" {
		e.FIELDS = "fields"
	}
	return []string{e.FIELDS}
}

func (e *Edit) RunAction(db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	return e.RunActionContext(context.Background(), db, t, args, extra...)
}

func (e *Edit) RunActionContext(ctx context.Context, db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	e.setDefaultElementNames()
	sql, labels := t.filterPars(args, e.FIELDS, e.getAllowed())

	ids := t.getIDVal(args, extra...)
	if !hasValue(ids) {
		return nil, errorMissingPk(t.TableName)
	}

	newExtra := t.byConstraint(args, extra...)
	where, extraValues := t.singleCondition(ids, t.TableName, newExtra)
	if where != "" {
		sql += "\nWHERE " + where
	}
	if t.dbDriver == Postgres {
		sql = questionMarkerNumber(sql)
	}

	return getSQL(ctx, db, t.logger, sql, labels, extraValues...)
}
