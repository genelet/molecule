package molecule

import (
	"context"
	"database/sql"
)

type Delete struct {
	Action
}

func (self *Delete) RunAction(db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Delete) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	ids := t.getIdVal(ARGS)
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
	dbi := &DBI{DB: db}
	if t.dbDriver == Postgres { sql = questionMarkerNumber(sql) }
	return nil, dbi.DoSQLContext(ctx, sql, values...)
}
