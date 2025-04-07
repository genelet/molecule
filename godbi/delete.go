package godbi

import (
	"context"
	"database/sql"
)

// Delete struct for row deletion by primary key
type Delete struct {
	Action
}

var _ Capability = (*Delete)(nil)

func (self *Delete) RunAction(db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Delete) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	ids := t.getIDVal(ARGS)
	if !hasValue(ids) {
		return nil, errorMissingPk(t.TableName)
	}

	sql := "DELETE FROM " + t.TableName
	where, values := t.singleCondition(ids, "", extra...)
	if where != "" {
		sql += "\nWHERE " + where
	} else {
		return nil, errorDeleteWhole(t.TableName)
	}
	dbi := &DBI{DB: db, logger: t.logger}
	if t.dbDriver == Postgres {
		sql = questionMarkerNumber(sql)
	}
	return nil, dbi.DoSQLContext(ctx, sql, values...)
}
