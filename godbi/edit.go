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

func (self *Edit) setDefaultElementNames() []string {
	if self.FIELDS == "" {
		self.FIELDS = "fields"
	}
	return []string{self.FIELDS}
}

func (self *Edit) RunAction(db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Edit) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	self.setDefaultElementNames()
	sql, labels := t.filterPars(ARGS, self.FIELDS, self.getAllowed())

	ids := t.getIDVal(ARGS, extra...)
	if !hasValue(ids) {
		return nil, errorMissingPk(t.TableName)
	}

	newExtra := t.byConstraint(ARGS, extra...)
	where, extraValues := t.singleCondition(ids, t.TableName, newExtra)
	if where != "" {
		sql += "\nWHERE " + where
	}
	if t.dbDriver == Postgres {
		sql = questionMarkerNumber(sql)
	}

	return getSQL(ctx, db, t.logger, sql, labels, extraValues...)
}
