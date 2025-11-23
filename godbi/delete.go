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

func (d *Delete) RunAction(db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	return d.RunActionContext(context.Background(), db, t, args, extra...)
}

func (d *Delete) RunActionContext(ctx context.Context, db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	ids := t.getIDVal(args)
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
	_, err := dbi.DoSQLContext(ctx, sql, values...)
	return nil, err
}
