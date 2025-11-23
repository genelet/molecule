package rdb

import (
	"database/sql"
	"strings"

	"github.com/genelet/molecule/godbi"
)

type postgres struct {
	database
}

func newPostgres(databaseName string) *postgres {
	team := database{DBDriver: godbi.Postgres, DatabaseName: databaseName}
	postgres := &postgres{team}
	postgres.database.schema = postgres
	return postgres
}

func postgresToNative(u string) string {
	arr := strings.Split(u, "(")
	key := strings.TrimSpace(arr[0])
	pre := map[string]string{
		"boolean":                     "bool",
		"bytea":                       "string",
		"character varying":           "string",
		"character":                   "string",
		"text":                        "string",
		"name":                        "string",
		"smallint":                    "int16",
		"integer":                     "int",
		"bigint":                      "int64",
		"real":                        "float32",
		"double precision":            "float64",
		"numeric":                     "float64",
		"decimal":                     "float64",
		"date":                        "string",
		"timestamp without time zone": "string",
		"timestamp with time zone":    "string",
		"time without time zone":      "string",
		"time with time zone":         "string",
		"json":                        "string",
		"jsonb":                       "string",
		"uuid":                        "string",
		"xml":                         "string",
	}
	if v, ok := pre[key]; ok {
		return v
	}
	return "string"
}

func (p *postgres) getTable(db *sql.DB, tableName string) (*godbi.Table, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.Select(&lists,
		`SELECT
	pg_catalog.quote_ident(a.attname) AS "COLUMN_NAME"
	, pg_catalog.format_type(a.atttypid, NULL) AS "TYPE_NAME"
	, pg_catalog.pg_get_expr(af.adbin, af.adrelid) AS "COLUMN_DEF"
	, CASE a.attnotnull WHEN 't' THEN 'NO' ELSE 'YES' END AS "IS_NULLABLE"
FROM pg_catalog.pg_type t
JOIN pg_catalog.pg_attribute a ON (t.oid = a.atttypid)
JOIN pg_catalog.pg_class c ON (a.attrelid = c.oid)
LEFT JOIN pg_catalog.pg_attrdef af ON (a.attnum = af.adnum AND a.attrelid = af.adrelid)
JOIN pg_catalog.pg_namespace n ON (n.oid = c.relnamespace)
WHERE a.attnum >= 0 AND c.relkind IN ('r','p','v','m','f')
AND pg_catalog.quote_ident(n.nspname) = 'public'
AND pg_catalog.quote_ident(pg_catalog.current_database()) = ?
AND pg_catalog.quote_ident(c.relname) = ?`, p.DatabaseName, tableName)
	if err != nil {
		return nil, err
	}

	var idauto string
	ref := make(map[string]*godbi.Col)

	for _, iitem := range lists {
		item := iitem.(map[string]any)
		field := toString(item["COLUMN_NAME"])
		col := &godbi.Col{
			ColumnName: field,
			Label:      field,
			TypeName:   postgresToNative(toString(item["TYPE_NAME"]))}

		if toString(item["IS_NULLABLE"]) == "NO" {
			col.Notnull = true
		}

		definition := toString(item["COLUMN_DEF"])
		if len(definition) > 7 && definition[0:7] == "nextval" {
			idauto = field
			col.Auto = true
			col.Notnull = true
		}
		if definition == "CURRENT_TIMESTAMP" {
			col.Auto = true
			col.Notnull = true
		}
		ref[field] = col
	}

	lists = make([]any, 0)
	err = dbi.Select(&lists,
		`SELECT kcu.column_name
FROM information_schema.table_constraints tco
JOIN information_schema.key_column_usage kcu
ON	kcu.constraint_schema = tco.constraint_schema AND
	kcu.constraint_name = tco.constraint_name
WHERE tco.constraint_type = 'PRIMARY KEY'
AND kcu.table_schema = 'public'
AND tco.constraint_catalog = ?
AND kcu.table_name = ?
ORDER BY kcu.ordinal_position`, p.DatabaseName, tableName)
	if err != nil {
		return nil, err
	}

	var pks []string
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		field := toString(item["column_name"])
		col := ref[field]
		col.Notnull = true
		col.Constraint = true
		ref[field] = col
		pks = append(pks, field)
	}

	var cols []*godbi.Col
	for _, col := range ref {
		cols = append(cols, col)
	}

	return &godbi.Table{
		TableName: tableName,
		Columns:   cols,
		Pks:       pks,
		IDAuto:    idauto}, nil
}

func (p *postgres) getFks(db *sql.DB, tableName string) ([]*godbi.Fk, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	// from child table: constraint_name | table_name  | column_name | foreign_table_name | foreign_column_name
	// poll_choice_poll_id..poll_poll_id | poll_choice | poll_id     | poll_question      | poll_id
	err := dbi.Select(&lists,
		`SELECT tc.constraint_name, tc.table_name, 
    kcu.column_name, 
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
	ON tc.constraint_name = kcu.constraint_name
	AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
	ON ccu.constraint_name = tc.constraint_name
	AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
AND tc.constraint_catalog = ?
AND tc.table_schema='public'
AND ccu.table_schema='public'
AND ccu.table_name=?`, p.DatabaseName, tableName)
	if err != nil {
		return nil, err
	}

	var fks []*godbi.Fk
	for _, iitem := range lists {
		item := iitem.(map[string]any)
		fkTable := item["foreign_table_name"].(string)
		fkColumn := item["foreign_column_name"].(string)
		column := item["column_name"].(string)
		if fkTable == tableName && fkColumn == column {
			continue
		}
		fks = append(fks, &godbi.Fk{FkTable: fkTable, FkColumn: fkColumn, Column: column})
	}

	return fks, nil
}

func (p *postgres) tableNames(db *sql.DB) ([]string, error) {
	dbi := &godbi.DBI{DB: db}
	lists := make([]any, 0)
	err := dbi.Select(&lists,
		`SELECT table_name FROM information_schema.tables
WHERE table_schema='public'
AND table_type='BASE TABLE'
AND table_catalog=?`, p.DatabaseName)
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

/*
psql
\l
\c tustin
SELECT c.conname                                 AS constraint_name,
   c.contype                                     AS constraint_type,
   sch.nspname                                   AS "self_schema",
   tbl.relname                                   AS "self_table",
   ARRAY_AGG(col.attname ORDER BY u.attposition) AS "self_columns",
   f_sch.nspname                                 AS "foreign_schema",
   f_tbl.relname                                 AS "foreign_table",
   ARRAY_AGG(f_col.attname ORDER BY f_u.attposition) AS "foreign_columns",
   pg_get_constraintdef(c.oid)                   AS definition
FROM pg_constraint c
       LEFT JOIN LATERAL UNNEST(c.conkey) WITH ORDINALITY AS u(attnum, attposition) ON TRUE
       LEFT JOIN LATERAL UNNEST(c.confkey) WITH ORDINALITY AS f_u(attnum, attposition) ON f_u.attposition = u.attposition
       JOIN pg_class tbl ON tbl.oid = c.conrelid
       JOIN pg_namespace sch ON sch.oid = tbl.relnamespace
       LEFT JOIN pg_attribute col ON (col.attrelid = tbl.oid AND col.attnum = u.attnum)
       LEFT JOIN pg_class f_tbl ON f_tbl.oid = c.confrelid
       LEFT JOIN pg_namespace f_sch ON f_sch.oid = f_tbl.relnamespace
       LEFT JOIN pg_attribute f_col ON (f_col.attrelid = f_tbl.oid AND f_col.attnum = f_u.attnum)
GROUP BY constraint_name, constraint_type, "self_schema", "self_table", definition, "foreign_schema", "foreign_table"
ORDER BY "self_schema", "self_table";
*/
