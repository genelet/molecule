package rdb

import (
	"database/sql"
	"strings"

	"github.com/genelet/molecule/godbi"
)

type mySQL struct {
	database
}

func newMySQL(databaseName string) *mySQL {
	team := database{DBDriver: godbi.MySQL, DatabaseName: databaseName}
	mysql := &mySQL{team}
	mysql.database.schema = mysql
	return mysql
}

func mysqlToNative(u string) string {
	arr := strings.Split(u, " ")
	ar2 := strings.Split(arr[0], "(")
	pre := map[string]string{
		"VARCHAR":  "string",
		"CHAR":     "string",
		"TEXT":     "string",
		"ENUM":     "string",
		"DOUBLE":   "float64",
		"FLOAT":    "float64",
		"BOOL":     "bool",
		"DATETIME": "string",
		"TINYINT":  "int8",
		"SMALLINT": "int16",
		"INT":      "int",
		"BIGINT":   "int64",
	}
	v, ok := pre[strings.ToUpper(strings.TrimSpace(ar2[0]))]
	if ok {
		return v
	}
	return ""
}

func (self *mySQL) getTable(db *sql.DB, tableName string) (*godbi.Table, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.SelectSQL(&lists,
		`DESC `+tableName,
		[]any{"Field", "Type", "Null", "Key", "Default", "Extra"})
	if err != nil {
		return nil, err
	}

	var pks []string
	var idauto string
	var cols []*godbi.Col

	for _, iitem := range lists {
		item := iitem.(map[string]any)
		field := item["Field"].(string)
		col := &godbi.Col{
			ColumnName: field,
			Label:      field,
			TypeName:   mysqlToNative(item["Type"].(string))}
		if item["Null"].(string) == "NO" {
			col.Notnull = true
		}
		if item["Key"].(string) == "PRI" {
			pks = append(pks, field)
			col.Notnull = true
			col.Constraint = true
		}
		if item["Extra"].(string) == "auto_increment" {
			idauto = field
			col.Auto = true
			col.Notnull = true
		} else if item["Default"] != nil && item["Default"].(string) == "CURRENT_TIMESTAMP" {
			col.Auto = true
			col.Notnull = true
		}
		cols = append(cols, col)
	}

	return &godbi.Table{
		TableName: tableName,
		Columns:   cols,
		Pks:       pks,
		IDAuto:    idauto}, nil
}

func (self *mySQL) getFks(db *sql.DB, tableName string) ([]*godbi.Fk, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.Select(&lists,
		`SELECT A.REFERENCED_TABLE_SCHEMA AS FKTABLE_SCHEM,
	A.REFERENCED_TABLE_NAME AS FKTABLE_NAME,
	A.REFERENCED_COLUMN_NAME AS FKCOLUMN_NAME,
	A.TABLE_SCHEMA AS PKTABLE_SCHEM,
	A.TABLE_NAME AS PKTABLE_NAME,
	A.COLUMN_NAME AS PKCOLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE A,
	INFORMATION_SCHEMA.TABLE_CONSTRAINTS B
WHERE (A.TABLE_SCHEMA = B.TABLE_SCHEMA)
AND (A.TABLE_NAME = B.TABLE_NAME)
AND (A.CONSTRAINT_NAME = B.CONSTRAINT_NAME)
AND (B.CONSTRAINT_TYPE IS NOT NULL)
AND A.REFERENCED_TABLE_NAME IS NOT NULL
AND A.TABLE_SCHEMA=?
AND A.TABLE_NAME=?`, self.DatabaseName, tableName)
	if err != nil {
		return nil, err
	}

	var fks []*godbi.Fk
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		fkTable := item["FKTABLE_NAME"].(string)
		fkColumn := item["FKCOLUMN_NAME"].(string)
		column := item["PKCOLUMN_NAME"].(string)
		if fkTable == tableName && fkColumn == column {
			continue
		}
		fks = append(fks, &godbi.Fk{FkTable: fkTable, FkColumn: fkColumn, Column: column})
	}

	return fks, nil
}

func (self *mySQL) getFwks(db *sql.DB, tableName string) ([]*godbi.Fk, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.Select(&lists,
		`SELECT A.REFERENCED_TABLE_SCHEMA AS PKTABLE_SCHEM,
	A.REFERENCED_TABLE_NAME AS PKTABLE_NAME,
	A.REFERENCED_COLUMN_NAME AS PKCOLUMN_NAME,
	A.TABLE_SCHEMA AS FKTABLE_SCHEM,
	A.TABLE_NAME AS FKTABLE_NAME,
	A.COLUMN_NAME AS FKCOLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE A,
	INFORMATION_SCHEMA.TABLE_CONSTRAINTS B
WHERE (A.TABLE_SCHEMA = B.TABLE_SCHEMA)
AND (A.TABLE_NAME = B.TABLE_NAME)
AND (A.CONSTRAINT_NAME = B.CONSTRAINT_NAME)
AND (B.CONSTRAINT_TYPE IS NOT NULL)
AND A.REFERENCED_TABLE_SCHEMA=?
AND A.REFERENCED_TABLE_NAME=?`, self.DatabaseName, tableName)
	if err != nil {
		return nil, err
	}

	var fks []*godbi.Fk
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		fkTable := item["FKTABLE_NAME"].(string)
		fkColumn := item["FKCOLUMN_NAME"].(string)
		column := item["PKCOLUMN_NAME"].(string)
		if fkTable == tableName && fkColumn == column {
			continue
		}
		fks = append(fks, &godbi.Fk{FkTable: fkTable, FkColumn: fkColumn, Column: column})
	}

	return fks, nil
}

func (self *mySQL) tableNames(db *sql.DB) ([]string, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.Select(&lists,
		`SELECT table_name AS table_name
FROM information_schema.tables
WHERE table_type='BASE TABLE' AND table_schema = ?`, self.DatabaseName)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		names = append(names, item["table_name"].(string))
	}
	return names, nil
}
