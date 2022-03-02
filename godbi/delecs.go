package godbi

import (
	"context"
	"database/sql"
	"strings"
)

type Delecs struct {
	Action
}

func (self *Delecs) RunAction(db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

// RunActionContext of Delecs is to populate Fks before making delete,
// so that Fks could be passed to other Prepared delete actions.
// If there is no Fks, we putput the input, then Delecs does nothing.
//
func (self *Delecs) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	dbi := &DBI{DB: db}
	lists := make([]interface{}, 0)
	if t.Fks == nil {
		return []interface{}{ARGS}, nil
	}
	str := ""
	var values []interface{}
	for _, pk := range t.Pks {
		if v, ok := ARGS[pk]; ok {
			if str != "" { str += ", " }
			str += pk + "=?"
			values = append(values, v)
		}
	}
	for _, fk := range t.Fks {
		name := fk.Column
		if v, ok := ARGS[name]; ok {
			if str != "" { str += ", " }
			str += name + "=?"
			values = append(values, v)
		}
	}
	if !hasValue(values) {
		return nil, errorMissingKeys(t.TableName)
	}
	if t.dbDriver == Postgres { str = questionMarkerNumber(str) }
	// should we add labels as well, if it != column name ?
	err := dbi.SelectContext(ctx, &lists, `SELECT ` + strings.Join(t.getKeyColumns(), ", ") + ` FROM ` + t.TableName + ` WHERE ` + str, values...)
	if err != nil { return nil, err }
	if hasValue(lists) && hasValue(lists[0]) {
		item := lists[0].(map[string]interface{})
		for k, v := range item {
			if v == nil {
				delete(item, k)
			}
		}
	}
	return lists, nil
}
