package rdb

import (
	"database/sql"
	"strings"

	"github.com/genelet/molecule/godbi"
)

type sQLite struct {
	database
}

func newSQLite(databaseName string) *sQLite {
	team := database{DBDriver: godbi.SQLite, DatabaseName: databaseName}
	sqlite := &sQLite{team}
	sqlite.database.schema = sqlite
	return sqlite
}

func sqliteToNative(u string) string {
	arr := strings.Split(u, " ")
	ar2 := strings.Split(arr[0], "(")
	pre := map[string]string{
		"INT":       "int",
		"INTEGER":   "int",
		"TINYINT":   "int8",
		"SMALLINT":  "int16",
		"MEDIUMINT": "int24",
		"BIGINT":    "int64",
		"UNSIGNED":  "int64",
		"INT2":      "bool",
		"INT8":      "int8",

		"CHARACTER": "string",
		"VARCHAR":   "string",
		"VARYING":   "string",
		"NCHAR":     "string",
		"NATIVE":    "string",
		"NVARCHAR":  "string",
		"TEXT":      "string",
		"CLOB":      "string",

		"BLOB": "string",

		"REAL":   "float64",
		"DOUBLE": "float64",
		"FLOAT":  "float32",

		"NUMERIC": "string",
		"DECIMAL": "string",
		"BOOLEAN": "bool",
		"DATE":    "string",
	}
	v, ok := pre[strings.ToUpper(strings.TrimSpace(ar2[0]))]
	if ok {
		return v
	}
	return ""
}

// .header on
// .mode column
// PRAGMA table_info(table_name);
func (s *sQLite) getTable(db *sql.DB, tableName string) (*godbi.Table, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.SelectSQL(&lists,
		`PRAGMA table_info(`+tableName+`)`,
		[]any{[2]string{"cid", "int"}, [2]string{"name", "string"}, "type", [2]string{"notnull", "int"}, [2]string{"default", "string"}, [2]string{"pk", "int"}})
	if err != nil {
		return nil, err
	}

	var pks []string
	idauto := "rowid"
	var cols []*godbi.Col

	for _, iitem := range lists {
		item := iitem.(map[string]any)
		field := toString(item["name"])
		col := &godbi.Col{
			ColumnName: field,
			Label:      field,
			TypeName:   sqliteToNative(toString(item["type"]))}
		if toString(item["default"]) == "CURRENT_TIMESTAMP" {
			col.Auto = true
		}
		pk := item["pk"].(int)
		notnull := item["notnull"].(int)
		if notnull == 1 {
			col.Notnull = true
			col.Constraint = true
		}
		if pk == 1 {
			pks = append(pks, field)
		}
		//if pk == 1 &&  notnull == 0 {
		//	idauto = field
		//}
		cols = append(cols, col)
	}
	return &godbi.Table{
		TableName: tableName,
		Columns:   cols,
		Pks:       pks,
		IDAuto:    idauto}, nil
}

// PRAGMA foreign_keys; 0 means disabled
// PRAGMA foreign_keys = ON;
// PRAGMA foreign_key_list(tableName);

func (s *sQLite) getFks(db *sql.DB, tableName string) ([]*godbi.Fk, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.SelectSQL(&lists,
		/*
		   sqlite> PRAGMA foreign_key_list(nodes);
		   id          seq         table       from        to          on_update   on_delete   match
		   ----------  ----------  ----------  ----------  ----------  ----------  --------    ------------

		   id and seq - each row the pragma returns refers to a single column, but foreign keys can be composite ones with multiple columns. id refers to the Nth key, and seq to the Nth column in that key. So a table with a single FK on two columns would have rows with id = 0 seq = 0 and id = 0 seq = 1.
		   table and to - The parent table and column in it the FK refers to. types.id in this case.
		   from the column in the child table that has the FK constraint. nodes.typeid in this case.
		   on_update and on_delete - the actions taken when the referenced foreign key is updated or deleted.
		   match - A SQL92 feature for foreign key actions related to null values that sqlite doesn't implement but does accept syntax-wise in case it ever does add support.
		*/
		`SELECT m.name, p."table", p."from", p."to"
FROM sqlite_master m
JOIN pragma_foreign_key_list(m.name) p ON m.name != p."table"
WHERE m.type = 'table'
AND m.name=?`, []any{"name", "table", "from", "to"}, tableName)
	if err != nil {
		return nil, err
	}

	var fks []*godbi.Fk
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		fkTable := toString(item["table"])
		fkColumn := toString(item["to"])
		column := toString(item["from"])
		if fkTable == tableName && fkColumn == column {
			continue
		}
		fks = append(fks, &godbi.Fk{FkTable: fkTable, FkColumn: fkColumn, Column: column})
	}

	return fks, nil
}

func (s *sQLite) tableNames(db *sql.DB) ([]string, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.Select(&lists,
		`SELECT name
FROM sqlite_master 
WHERE type ='table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		names = append(names, item["name"].(string))
	}
	return names, nil
}
