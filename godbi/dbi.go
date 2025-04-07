// Package godbi provides a generic database interface
package godbi

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// DBI embeds GO's generic SQL handler and
// adds few functions for database executions and queries.
type DBI struct {
	// Embedding the generic database handle.
	*sql.DB
	// LastID: the last auto id inserted, if the database provides
	LastID int64
	// Affected: the number of rows affected by the last SQL execution
	Affected int64
	// Slogger: a logger for SQL execution
	logger Slogger
}

// TxSQL is the same as DoSQL, but use transaction
func (self *DBI) TxSQL(query string, args ...any) error {
	return self.TxSQLContext(context.Background(), query, args...)
}

// TxSQLContext is the same as DoSQLContext, but use transaction
func (self *DBI) TxSQLContext(ctx context.Context, query string, args ...any) error {
	tx, err := self.DB.Begin()
	if err != nil {
		return err
	}
	//defer tx.Rollback()

	if self.logger != nil {
		self.logger.Debug("godbi.DBI", "SQL", query, "ARGS", args)
	}
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errorRollback(err, rollbackErr)
		} else {
			return err
		}
	}

	// commit the trasaction
	if err := tx.Commit(); err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	self.LastID = lastID

	affected, err := res.RowsAffected()
	if err == nil {
		self.Affected = affected
	}

	return err
}

// InsertSerial insert a SQL row into Postgres table with Serail , only save the last inserted ID
func (self *DBI) InsertSerial(query string, args ...any) error {
	return self.InsertSerialContext(context.Background(), query, args...)
}

// InsertSerialContext insert a SQL into Postgres table with Serial, only save the last inserted ID
func (self *DBI) InsertSerialContext(ctx context.Context, query string, args ...any) error {
	stmt, err := self.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var lastID int64
	if self.logger != nil {
		self.logger.Debug("godbi.DBI", "SQL", query, "ARGS", args)
	}
	err = stmt.QueryRowContext(ctx, args...).Scan(&lastID)
	if err != nil {
		return err
	}
	self.LastID = lastID

	return nil
}

// InsertID executes a SQL the same as DB's Exec, only save the last inserted ID
func (self *DBI) InsertID(query string, args ...any) error {
	return self.InsertIDContext(context.Background(), query, args...)
}

// InsertIDContext executes a SQL the same as DB's Exec, only save the last inserted ID
func (self *DBI) InsertIDContext(ctx context.Context, query string, args ...any) error {
	if self.logger != nil {
		self.logger.Debug("godbi.DBI", "SQL", query, "ARGS", args)
	}
	res, err := self.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	self.LastID = lastID

	affected, err := res.RowsAffected()
	if err == nil {
		self.Affected = affected
	}

	return nil
}

// DoSQL executes a SQL the same as DB's Exec, only save the last inserted ID
func (self *DBI) DoSQL(query string, args ...any) error {
	return self.DoSQLContext(context.Background(), query, args...)
}

// DoSQLContext executes a SQL the same as DB's Exec, only save the last inserted ID
func (self *DBI) DoSQLContext(ctx context.Context, query string, args ...any) error {
	if self.logger != nil {
		self.logger.Debug("godbi.DBI", "SQL", query, "ARGS", args)
	}
	res, err := self.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	self.LastID = lastID
	affected, err := res.RowsAffected()
	if err == nil {
		self.Affected = affected
	}

	return err
}

// DoSQLs executes multiple rows using the same prepared statement,
// The rows are expressed as a slice of slices.
func (self *DBI) DoSQLs(query string, args ...[]any) error {
	return self.DoSQLsContext(context.Background(), query, args...)
}

// DoSQLsContext executes multiple rows using the same prepared statement,
// The rows are expressed as a slice of slices.
func (self *DBI) DoSQLsContext(ctx context.Context, query string, args ...[]any) error {
	n := len(args)
	if n == 0 {
		return self.DoSQLContext(ctx, query)
	} else if n == 1 {
		return self.DoSQLContext(ctx, query, args[0]...)
	}

	sth, err := self.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	var res sql.Result
	for _, once := range args {
		res, err = sth.ExecContext(ctx, once...)
		if err != nil {
			return err
		}
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	self.LastID = lastID
	affected, err := res.RowsAffected()
	if err == nil {
		self.Affected = affected
	}
	if err != nil {
		return err
	}

	sth.Close()
	return nil
}

// Select returns queried data 'list' as a slice of maps,
// whose data types are determined dynamically by the generic handle.
func (self *DBI) Select(lists *[]any, query string, args ...any) error {
	return self.SelectContext(context.Background(), lists, query, args...)
}

// SelectContext returns queried data 'lists' as a slice of maps,
// whose data types are determined dynamically by the generic handle.
func (self *DBI) SelectContext(ctx context.Context, lists *[]any, query string, args ...any) error {
	return self.SelectSQLContext(ctx, lists, query, nil, args...)
}

func getLabels(labels []any) ([]string, []string) {
	if len(labels) == 0 {
		return nil, nil
	}

	var selectLabels []string
	var typeLabels []string
	for _, vs := range labels {
		switch v := vs.(type) {
		case []any:
			selectLabels = append(selectLabels, v[0].(string))
			if len(v) > 1 {
				typeLabels = append(typeLabels, v[1].(string))
			} else {
				typeLabels = append(typeLabels, "")
			}
		case [2]string:
			selectLabels = append(selectLabels, v[0])
			if len(v) > 1 {
				typeLabels = append(typeLabels, v[1])
			} else {
				typeLabels = append(typeLabels, "")
			}
		default:
			selectLabels = append(selectLabels, vs.(string))
			typeLabels = append(typeLabels, "")
		}
	}

	return selectLabels, typeLabels
}

// SelectSQL returns queried data 'list' as a slice of maps.
// The map keys and their data types are pre-defined in 'labels',
// expressed as a slice of interfaces:
//  1. when an interface is a string, this is the key name.
//     The data types are determined dynamically by the generic handler.
//  2. when an interface is a 2-string slice, the first element is the key
//     and the second the data type in "int64", "int", "string" etc.
func (self *DBI) SelectSQL(lists *[]any, query string, labels []any, args ...any) error {
	return self.SelectSQLContext(context.Background(), lists, query, labels, args...)
}

// SelectSQLContext returns queried data 'list' as a slice of maps.
// The map keys and their data types are pre-defined in 'labels',
// expressed as a slice of interfaces:
//  1. when an interface is a string, this is the key name.
//     The data types are determined dynamically by the generic handler.
//  2. when an interface is a 2-string slice, the first element is the key
//     and the second the data type in "int64", "int", "string" etc.
func (self *DBI) SelectSQLContext(ctx context.Context, lists *[]any, query string, labels []any, args ...any) error {
	if self.logger != nil {
		self.logger.Debug("godbi.DBI", "SQL", query, "ARGS", args)
	}
	rows, err := self.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return self.pickup(rows, lists, labels, query)
}

func (self *DBI) pickup(rows *sql.Rows, lists *[]any, labels []any, query string) error {
	selectLabels, typeLabels := getLabels(labels)

	var err error
	if selectLabels == nil {
		if selectLabels, err = rows.Columns(); err != nil {
			return err
		}
		typeLabels = make([]string, len(selectLabels))
	}

	names := make([]any, len(selectLabels))
	x := make([]any, len(selectLabels))
	for i := range selectLabels {
		switch typeLabels[i] {
		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "int64":
			x[i] = new(sql.NullInt64)
		case "float32", "float64":
			x[i] = new(sql.NullFloat64)
		case "bool":
			x[i] = new(sql.NullBool)
		case "time":
			x[i] = new(sql.NullTime)
		case "string", "[]byte":
			x[i] = new(sql.NullString)
		default:
			x[i] = &names[i]
		}
	}

	for rows.Next() {
		if err = rows.Scan(x...); err != nil {
			return err
		}
		res := make(map[string]any)
		for j, v := range selectLabels {
			switch typeLabels[j] {
			case "int":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = int(x.Int64)
				}
			case "int8":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = int8(x.Int64)
				}
			case "int16":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = int16(x.Int64)
				}
			case "int32":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = int32(x.Int64)
				}
			case "uint":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = uint(x.Int64)
				}
			case "uint8":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = uint8(x.Int64)
				}
			case "uint16":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = uint16(x.Int64)
				}
			case "uint32":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = uint32(x.Int64)
				}
			case "int64":
				x := x[j].(*sql.NullInt64)
				if x.Valid {
					res[v] = x.Int64
				}
			case "float32":
				x := x[j].(*sql.NullFloat64)
				if x.Valid {
					res[v] = float32(x.Float64)
				}
			case "float64":
				x := x[j].(*sql.NullFloat64)
				if x.Valid {
					res[v] = x.Float64
				}
			case "bool":
				x := x[j].(*sql.NullBool)
				if x.Valid {
					res[v] = x.Bool
				}
			case "time":
				x := x[j].(*sql.NullTime)
				if x.Valid {
					res[v] = x.Time
				}
			case "string", "[]byte":
				x := x[j].(*sql.NullString)
				if x.Valid {
					res[v] = x.String
				}
			default:
				name := names[j]
				res[v] = name
				if name != nil {
					switch val := name.(type) {
					case []uint8:
						res[v] = string(val)
					case string:
						res[v] = val
					default:
						res[v] = fmt.Sprintf("%v", val)
					}
				}
			}
		}
		*lists = append(*lists, res)
	}
	rows.Close()
	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// GetSQL returns single row 'res'.
func (self *DBI) GetSQL(res map[string]any, query string, labels []any, args ...any) error {
	return self.GetSQLContext(context.Background(), res, query, labels, args...)
}

// GetSQLContext returns single row 'res'.
func (self *DBI) GetSQLContext(ctx context.Context, res map[string]any, query string, labels []any, args ...any) error {
	lists := make([]any, 0)
	if err := self.SelectSQLContext(ctx, &lists, query, labels, args...); err != nil {
		return err
	}
	if len(lists) >= 1 {
		for k, v := range lists[0].(map[string]any) {
			if v != nil {
				res[k] = v
			}
		}
	}
	return nil
}

func sp(procName string, labels []any, n int) (string, string) {
	names, _ := getLabels(labels)
	strQ := strings.Join(strings.Split(strings.Repeat("?", n), ""), ",")
	str := procName + "(" + strQ
	strN := "@" + strings.Join(names, ",@")
	if names != nil {
		str += ", " + strN
	}
	return str + ")", "SELECT " + strN
}

/*

// doProc runs the stored procedure 'procName' and outputs
// the OUT data as map whose keys are in 'names'.
func (self *DBI) doProc(res map[string]any, procName string, names []any, args ...any) error {
	return self.doProcContext(context.Background(), res, procName, names, args...)
}

// doProcContext runs the stored procedure 'procName' and outputs
// the OUT data as map whose keys are in 'names'.
func (self *DBI) doProcContext(ctx context.Context, res map[string]any, procName string, names []any, args ...any) error {
	str, strN := sp(procName, names, len(args))
	if err := self.DoSQLContext(ctx, str, args...); err != nil {
		return err
	}
	return self.GetSQLContext(ctx, res, strN, names)
}

// txProc runs the stored procedure 'procName' in transaction
// and outputs the OUT data as map whose keys are in 'names'.
func (self *DBI) txProc(res map[string]any, procName string, names []any, args ...any) error {
	return self.txProcContext(context.Background(), res, procName, names, args...)
}

// txProcContext runs the stored procedure 'procName' in transaction
// and outputs the OUT data as map whose keys are in 'names'.
func (self *DBI) txProcContext(ctx context.Context, res map[string]any, procName string, names []any, args ...any) error {
	str, strN := sp(procName, names, len(args))
	if err := self.TxSQLContext(ctx, str, args...); err != nil {
		return err
	}
	return self.GetSQLContext(ctx, res, strN, names)
}

// selectProc runs the stored procedure 'procName'.
// The result, 'lists', is received as slice of map whose key
// names and data types are defined in 'labels'.
func (self *DBI) selectProc(lists *[]any, procName string, labels []any, args ...any) error {
	return self.selectProcContext(context.Background(), lists, procName, labels, args...)
}

*/

// selectProcContext runs the stored procedure 'procName'.
// The result, 'lists', is received as slice of map whose key
// names and data types are defined in 'labels'.
func (self *DBI) selectProcContext(ctx context.Context, lists *[]any, procName string, labels []any, args ...any) error {
	return self.selectDoProcContext(ctx, lists, nil, nil, procName, labels, args...)
}

/*

// getProc returns single row from stored procedure into 'res'.
func (self *DBI) getProc(res map[string]any, procName string, labels []any, args ...any) error {
	return self.getProcContext(context.Background(), res, procName, labels, args...)
}

// getProcContext returns single row from stored procedure into 'res'.
func (self *DBI) getProcContext(ctx context.Context, res map[string]any, procName string, labels []any, args ...any) error {
	lists := make([]any, 0)
	if err := self.selectProcContext(ctx, &lists, procName, labels, args...); err != nil {
		return err
	}
	if len(lists) >= 1 {
		for k, v := range lists[0].(map[string]any) {
			if v != nil {
				res[k] = v
			}
		}
	}
	return nil
}

// selectDoProc runs the stored procedure 'procName'.
// The result, 'lists', is received as slice of map whose key names and data
// types are defined in 'labels'. The OUT data, 'hash', is received as map
// whose keys are in 'names'.
func (self *DBI) selectDoProc(lists *[]any, hash map[string]any, names []any, procName string, labels []any, args ...any) error {
	return self.selectDoProcContext(context.Background(), lists, hash, names, procName, labels, args...)
}

*/

// selectDoProcContext runs the stored procedure 'procName'.
// The result, 'lists', is received as slice of map whose key names and data
// types are defined in 'labels'. The OUT data, 'hash', is received as map
// whose keys are in 'names'.
func (self *DBI) selectDoProcContext(ctx context.Context, lists *[]any, hash map[string]any, names []any, procName string, labels []any, args ...any) error {
	str, strN := sp(procName, names, len(args))
	if err := self.SelectSQLContext(ctx, lists, str, labels, args...); err != nil {
		return err
	}
	if hash == nil {
		return nil
	}
	return self.GetSQLContext(ctx, hash, strN, names)
}
