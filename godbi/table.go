package godbi

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// Col defines table column in GO struct
type Col struct {
	ColumnName string `json:"columnName" hcl:"columnName,label"`
	TypeName   string `json:"typeName" hcl:"typeName,optional"`
	Label      string `json:"label" hcl:"columnLabel,optional"`
	Notnull    bool   `json:"notnull,omitempty" hcl:"notnull,optional"`
	Constraint bool   `json:"constraint,omitempty" hcl:"constraint,optional"`
	Auto       bool   `json:"auto,omitempty" hcl:"auto,optional"`
	// true for a one-to-may recurse column
	Recurse bool `json:"recurse,omitempty" hcl:"recurse,optional"`
}

// Fk defines foreign key struct
type Fk struct {
	// the parent table name
	FkTable string `json:"fkTable" hcl:"fkTable,optional"`
	// the parant table column
	FkColumn string `json:"fkColumn" hcl:"fkColumn,optional"`
	// column of this table
	Column string `json:"column" hcl:"column,optional"`
}

// Table defines a table by name, columns, primary key, foreign keys, auto id and unique columns
type Table struct {
	TableName string   `json:"tableName" hcl:"tableName,optional"`
	Columns   []*Col   `json:"columns" hcl:"columns,block"`
	Pks       []string `json:"pks,omitempty" hcl:"pks,optional"`
	IDAuto    string   `json:"idAuto,omitempty" hcl:"idAuto,optional"`
	Fks       []*Fk    `json:"fks,omitempty" hcl:"fks,block"`
	Uniques   []string `json:"uniques,omitempty" hcl:"uniques,optional"`
	dbDriver  DBType
	logger    Slogger
}

// SetLogger sets the logger
func (self *Table) SetLogger(logger Slogger) {
	self.logger = logger
}

// GetLogger gets the logger
func (self *Table) GetLogger() Slogger {
	return self.logger
}

// IsRecursive indicates if table references to itself in one to multiple relations
func (self *Table) IsRecursive() bool {
	for _, col := range self.Columns {
		if col.ColumnName == self.Pks[0] && col.Recurse {
			return true
		}
	}
	return false
}

// RecursiveColumn returns the name of the resursive column
func (self *Table) RecursiveColumn() string {
	for _, col := range self.Columns {
		if col.ColumnName == self.Pks[0] || !col.Recurse {
			continue
		}
		return col.ColumnName
	}
	return ""
}

// SetDBDriver sets the driver type
func (self *Table) SetDBDriver(driver DBType) {
	self.dbDriver = driver
}

func (self *Table) byConstraint(ARGS map[string]any, extra ...map[string]any) map[string]any {
	var output map[string]any
	for k, v := range ARGS {
		find := false
		for _, col := range self.Columns {
			if col.Label == k && col.Constraint {
				find = true
			}
		}
		if !find {
			continue
		}
		if output == nil {
			output = make(map[string]any)
		}
		output[k] = v
	}
	if hasValue(extra) && hasValue(extra[0]) {
		for k, v := range extra[0] {
			if output == nil {
				output = make(map[string]any)
			}
			output[k] = v
		}
	}
	return output
}

// refreshes args by checking if column's label exists a key.
// If it exists,
// if force is true, forcefully sets the column using label's value;
// if force is not set, optionally set the column.
func (self *Table) refreshArgs(args any, force ...bool) any {
	if args == nil {
		return args
	}

	cut := func(item map[string]any, force ...bool) map[string]any {
		newArgs := make(map[string]any)
		for k, v := range item {
			newArgs[k] = v
		}
		for _, col := range self.Columns {
			v, ok := item[col.Label]
			if !ok {
				continue
			}
			if force != nil && force[0] {
				newArgs[col.ColumnName] = v
			} else if _, ok := item[col.ColumnName]; !ok {
				newArgs[col.ColumnName] = v
			}
		}
		return newArgs
	}

	switch t := args.(type) {
	case []any:
		var lists any
		for _, item := range t {
			got := cut(item.(map[string]any), force...)
			if lists == nil {
				lists = []map[string]any{got}
			} else if !grepMap(lists.([]map[string]any), got) {
				lists = append(lists.([]map[string]any), got)
			}
		}
		return lists
	case []map[string]any:
		var lists any
		for _, item := range t {
			got := cut(item, force...)
			if lists == nil {
				lists = []map[string]any{got}
			} else if !grepMap(lists.([]map[string]any), got) {
				lists = append(lists.([]map[string]any), got)
			}
		}
		return lists
	case map[string]any:
		return cut(t, force...)
	default:
	}

	return nil
}

func (self *Table) getKeyColumns() []string {
	labels := make(map[string]bool)
	for _, pk := range self.Pks {
		labels[pk] = true
	}
	if self.IDAuto != "" {
		labels[self.IDAuto] = true
	}
	if self.Fks != nil {
		for _, fk := range self.Fks {
			labels[fk.Column] = true
		}
	}

	var outs []string
	for k := range labels {
		outs = append(outs, k)
	}
	return outs
}

func (self *Table) getFv(ARGS map[string]any, allowed map[string]bool) (map[string]any, bool) {
	fieldValues := make(map[string]any)
	for f, l := range self.insertCols(allowed) {
		v, ok := ARGS[f]
		if !ok {
			v, ok = ARGS[l]
		}
		if ok {
			switch t := v.(type) {
			case []map[string]any, map[string]any:
			case bool:
				switch self.dbDriver {
				case SQLite, TSNano:
					if t {
						fieldValues[f] = 1
					} else {
						fieldValues[f] = 0
					}
				default:
					if t {
						fieldValues[f] = "true"
					} else {
						fieldValues[f] = "false"
					}
				}
			default:
				fieldValues[f] = t
			}
		}
	}

	allAuto := true
	for _, col := range self.Columns {
		if !col.Auto && col.Notnull {
			allAuto = false
			break
		}
	}

	return fieldValues, allAuto
}

func (self *Table) checkNull(ARGS map[string]any, extra ...map[string]any) error {
	for _, col := range self.Columns {
		if !col.Notnull || col.Auto {
			continue
		} // the column is ok with null
		err := errorNoSuchColumn(col.ColumnName)
		if _, ok := ARGS[col.ColumnName]; !ok {
			if hasValue(extra) && hasValue(extra[0]) {
				if _, ok = extra[0][col.ColumnName]; !ok {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func (self *Table) insertCols(allowed map[string]bool) map[string]string {
	cols := make(map[string]string)
	for _, col := range self.Columns {
		if col.Auto {
			continue
		}
		if allowed != nil && !allowed[col.Label] {
			continue
		}
		cols[col.ColumnName] = col.Label
	}
	return cols
}

func (self *Table) insertHashContext(ctx context.Context, db *sql.DB, args map[string]any) (int64, error) {
	var fields []string
	var values []any
	if self.IDAuto != "" && self.dbDriver == TSNano {
		fields = append(fields, self.IDAuto)
		values = append(values, time.Now().UnixNano()/int64(time.Millisecond))
	}
	for k, v := range args {
		if v != nil {
			fields = append(fields, k)
			values = append(values, v)
		}
	}

	sql := "INSERT INTO " + self.TableName + " (" + strings.Join(fields, ", ") + ") VALUES (" + strings.Join(strings.Split(strings.Repeat("?", len(fields)), ""), ",") + ")"

	dbi := &DBI{DB: db, logger: self.logger}
	var err error
	switch self.dbDriver {
	case Postgres:
		sql = questionMarkerNumber(sql)
		if self.IDAuto != "" {
			sql += " RETURNING " + self.IDAuto
			err = dbi.InsertSerialContext(ctx, sql, values...)
		} else {
			err = dbi.DoSQLContext(ctx, sql, values...)
		}
	case SQLite:
		if hasValue(values) {
			err = dbi.InsertIDContext(ctx, sql, values...)
		} else {
			err = dbi.InsertIDContext(ctx, "INSERT INTO "+self.TableName+" DEFAULT VALUES")
		}
	case SQLRaw, TSNano:
		err = dbi.DoSQLContext(ctx, sql, values...)
	default:
		err = dbi.InsertIDContext(ctx, sql, values...)
	}
	if err != nil {
		return 0, err
	}
	return dbi.LastID, nil
}

func (self *Table) updateHashNullsContext(ctx context.Context, db *sql.DB, args map[string]any, ids []any, empties []string, extra ...map[string]any) error {
	if !hasValue(args) {
		return errorEmptyInput(self.TableName)
	}
	for _, k := range self.Pks {
		if grep(empties, k) {
			return errorMissingPk(self.TableName)
		}
	}

	n := len(args)
	fields := make([]string, n)
	field0 := make([]string, n)
	values := make([]any, n)
	i := 0
	for k, v := range args {
		fields[i] = k
		field0[i] = k + "=?"
		values[i] = v
		i++
	}

	sql := "UPDATE " + self.TableName + " SET " + strings.Join(field0, ", ")
	for _, v := range empties {
		if _, ok := args[v]; ok {
			continue
		}
		sql += ", " + v + "=NULL"
	}

	where, extraValues := self.singleCondition(ids, "", extra...)
	if where != "" {
		sql += "\nWHERE " + where
		values = append(values, extraValues...)
	}

	dbi := &DBI{DB: db, logger: self.logger}
	if self.dbDriver == Postgres {
		sql = questionMarkerNumber(sql)
	}
	return dbi.DoSQLContext(ctx, sql, values...)
}

func (self *Table) insupdTableContext(ctx context.Context, db *sql.DB, args map[string]any) (int64, error) {
	changed := int64(0)
	s := "SELECT " + strings.Join(self.Pks, ", ") + " FROM " + self.TableName + "\nWHERE "
	var v []any
	if self.Uniques == nil {
		return changed, errorNoUniqueKey(self.TableName)
	}
	for i, val := range self.Uniques {
		if i > 0 {
			s += " AND "
		}
		s += val + "=?"
		if x, ok := args[val]; ok {
			v = append(v, x)
		} else {
			return changed, errorEmptyInput(val)
		}
	}

	lists := make([]any, 0)
	dbi := &DBI{DB: db, logger: self.logger}
	if self.dbDriver == Postgres {
		s = questionMarkerNumber(s)
	}
	err := dbi.SelectContext(ctx, &lists, s, v...)
	if err != nil {
		return changed, err
	}
	if len(lists) > 1 {
		return changed, errorNotUnique(self.TableName)
	}

	if len(lists) == 1 {
		ids := make([]any, 0)
		for _, k := range self.Pks {
			ids = append(ids, lists[0].(map[string]any)[k])
		}
		err = self.updateHashNullsContext(ctx, db, args, ids, nil)
		if err == nil && self.IDAuto != "" {
			sql := "SELECT " + self.IDAuto + " FROM " + self.TableName + "\nWHERE " + strings.Join(self.Pks, "=? AND ") + "=?"
			if self.dbDriver == Postgres {
				sql = questionMarkerNumber(sql)
			}
			err = db.QueryRowContext(ctx, sql, ids...).Scan(&changed)
			return changed, err
		}
	} else {
		changed, err = self.insertHashContext(ctx, db, args)
	}

	return changed, err
}

func (self *Table) totalHashContext(ctx context.Context, db *sql.DB, v any, extra ...map[string]any) error {
	sql := "SELECT COUNT(*) FROM " + self.TableName

	if hasValue(extra) {
		where, values := selectCondition(extra[0], "")
		if where != "" {
			sql += "\nWHERE " + where
		}
		if self.dbDriver == Postgres {
			sql = questionMarkerNumber(sql)
		}
		return db.QueryRowContext(ctx, sql, values...).Scan(v)
	}

	return db.QueryRowContext(ctx, sql).Scan(v)
}

func (self *Table) getIDVal(ARGS map[string]any, extra ...map[string]any) []any {
	if hasValue(extra) {
		return properValues(self.Pks, ARGS, extra[0])
	}
	return properValues(self.Pks, ARGS, nil)
}

func (self *Table) singleCondition(ids []any, table string, extra ...map[string]any) (string, []any) {
	keys := self.Pks
	sql := ""
	var extraValues []any

	for i, item := range keys {
		val := ids[i]
		if i == 0 {
			sql = "("
		} else {
			sql += " AND "
		}
		switch idValues := val.(type) {
		case []any:
			n := len(idValues)
			sql += item + " IN (" + strings.Join(strings.Split(strings.Repeat("?", n), ""), ",") + ")"
			extraValues = append(extraValues, idValues...)
		default:
			sql += item + " =?"
			extraValues = append(extraValues, val)
		}
	}
	sql += ")"

	if hasValue(extra) && hasValue(extra[0]) {
		s, arr := selectCondition(extra[0], table)
		sql += " AND " + s
		extraValues = append(extraValues, arr...)
	}

	return sql, extraValues
}

func properValue(u string, ARGS map[string]any, extra map[string]any) any {
	if !hasValue(extra) {
		return ARGS[u]
	}
	if val, ok := extra[u]; ok {
		return val
	}
	return ARGS[u]
}

func properValues(us []string, ARGS map[string]any, extra map[string]any) []any {
	outs := make([]any, len(us))
	if !hasValue(extra) {
		for i, u := range us {
			outs[i] = ARGS[u]
		}
		return outs
	}
	for i, u := range us {
		if val, ok := extra[u]; ok {
			outs[i] = val
		} else {
			outs[i] = ARGS[u]
		}
	}
	return outs
}

func properValuesHash(vs []string, ARGS map[string]any, extra map[string]any) map[string]any {
	values := properValues(vs, ARGS, extra)
	hash := make(map[string]any)
	for i, v := range vs {
		hash[v] = values[i]
	}
	return hash
}

func selectCondition(extra map[string]any, table string) (string, []any) {
	sql := ""
	var values []any
	i := 0

	for field, valueInterface := range extra {
		if i > 0 {
			sql += " AND "
		}
		i++
		sql += "("

		if table != "" {
			if !strings.Contains(field, ".") {
				field = table + "." + field
			}
		}
		switch value := valueInterface.(type) {
		case []int:
			n := len(value)
			sql += field + " IN (" + strings.Join(strings.Split(strings.Repeat("?", n), ""), ",") + ")"
			for _, v := range value {
				values = append(values, v)
			}
		case []int64:
			n := len(value)
			sql += field + " IN (" + strings.Join(strings.Split(strings.Repeat("?", n), ""), ",") + ")"
			for _, v := range value {
				values = append(values, v)
			}
		case []string:
			n := len(value)
			sql += field + " IN (" + strings.Join(strings.Split(strings.Repeat("?", n), ""), ",") + ")"
			for _, v := range value {
				values = append(values, v)
			}
		case string:
			n := len(field)
			if n >= 5 && field[(n-5):] == "_gsql" {
				sql += value
			} else {
				sql += field + " =?"
				values = append(values, value)
			}
		default:
			sql += field + " =?"
			values = append(values, value)
		}
		sql += ")"
	}

	return sql, values
}

func (self *Table) filterPars(ARGS map[string]any, fieldsName string, allowed map[string]bool) (string, []any) {
	if allowed == nil {
		allowed = make(map[string]bool)
		for _, col := range self.Columns {
			allowed[col.Label] = true
		}
	}

	var fields map[string]bool
	if hasValue(ARGS) && hasValue(ARGS[fieldsName]) {
		fields = make(map[string]bool)
		for _, item := range strings.Split(ARGS[fieldsName].(string), ",") {
			if allowed[item] {
				fields[item] = true
			}
		}
	} else {
		fields = allowed
	}

	var keys []string
	var labels []any
	for _, col := range self.Columns {
		label := col.Label
		if fields == nil || fields[label] {
			keys = append(keys, col.ColumnName)
			labels = append(labels, [2]string{label, col.TypeName})
		}
	}

	return "SELECT " + strings.Join(keys, ", ") + "\nFROM " + self.TableName, labels
}

/*
func (self *Table) filterPars(ARGS map[string]any, fieldsName string, rest ...any) (string, []any, string) {
	var fields map[string]bool
	if v, ok := ARGS[fieldsName]; ok {
		fields := make(map[string]bool)
		for _, item := range strings.Split(v.(string), ",") {
			fields[item] = true
		}
	}

	var keys []string
	var labels []any
	if rest != nil && len(rest) == 2 {
		for k, label := range rest[1].(map[string]string) {
			if fields == nil || fields[label] == true {
				keys = append(keys, k)
				if len(rest) > 2 && rest[2] != nil {
					labels = append(labels, [2]string{label, rest[2].(map[string]string)[k]})
				} else {
					labels = append(labels, label)
				}
			}
		}
	} else {
		var allowed map[string]bool
		if rest != nil {
			allowed = rest[0].(map[string]bool)
		}
		for _, col := range self.Columns {
			label := col.Label
			if allowed != nil && !allowed[label] {
				continue
			}
			if fields == nil || fields[label] == true {
				keys = append(keys, col.ColumnName)
				labels = append(labels, [2]string{label, col.TypeName})
			}
		}
	}
	sql := strings.Join(keys, ", ")

	var table string
	if rest != nil && len(rest) == 2 {
		joints := rest[0].([]*Joint)
		sql = "SELECT " + sql + "\nFROM " + joinString(joints)
		table = joints[0].getAlias()
	} else {
		sql = "SELECT " + sql + "\nFROM " + self.TableName
	}

	return sql, labels, table
}
*/

func fromFv(fv map[string]any) []any {
	return []any{fv}
}
