package godbi

import (
	"context"
	"database/sql"
	"strings"
)

// Delecs is a special Topics action that returns all foreign keys
type Delecs struct {
	Action
}

var _ Capability = (*Delecs)(nil)

func (d *Delecs) RunAction(db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	return d.RunActionContext(context.Background(), db, t, args, extra...)
}

// RunActionContext of Delecs is to populate Fks before making delete,
// so that Fks could be passed to other Prepared delete actions.
// If there is no Fks, we putput the input, then Delecs does nothing.
func (d *Delecs) RunActionContext(ctx context.Context, db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	dbi := &DBI{DB: db, logger: t.logger}
	lists := make([]any, 0)
	if t.Fks == nil {
		return []any{args}, nil
	}
	str := ""
	var values []any
	for _, pk := range t.Pks {
		if v, ok := args[pk]; ok {
			if str != "" {
				str += ", "
			}
			str += pk + "=?"
			values = append(values, v)
		}
	}
	for _, fk := range t.Fks {
		name := fk.Column
		if v, ok := args[name]; ok {
			if str != "" {
				str += ", "
			}
			str += name + "=?"
			values = append(values, v)
		}
	}
	if !hasValue(values) {
		return nil, errorMissingKeys(t.TableName)
	}
	if t.dbDriver == Postgres {
		str = questionMarkerNumber(str)
	}
	// should we add labels as well, if it != column name ?
	err := dbi.SelectContext(ctx, &lists, `SELECT `+strings.Join(t.getKeyColumns(), ", ")+` FROM `+t.TableName+` WHERE `+str, values...)
	if err != nil {
		return nil, err
	}
	if hasValue(lists) && hasValue(lists[0]) {
		item := lists[0].(map[string]any)
		for k, v := range item {
			if v == nil {
				delete(item, k)
			}
		}
	}
	return lists, nil
}
