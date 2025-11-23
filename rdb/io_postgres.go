package rdb

import (
	"database/sql"
	"io"
	"strings"

	"github.com/genelet/molecule/godbi"

	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/dialect"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/genelet/sqlproto/ast"
	"github.com/genelet/sqlproto/light"
	"github.com/genelet/sqlproto/xlight"
)

type postgresIO struct {
	database
	tables map[string]*godbi.Table
	fks    map[string][]*godbi.Fk
}

func newPostgresIO(databaseName string, src io.Reader) (*postgresIO, error) {
	parser, err := xsqlparser.NewParser(src, &dialect.PostgresqlDialect{}, xsqlparser.ParseComment())
	if err != nil {
		return nil, err
	}
	stmts, err := parser.ParseSQL()
	if err != nil {
		return nil, err
	}

	tables := make(map[string]*godbi.Table)
	fks := make(map[string][]*godbi.Fk)

	for _, stmt := range stmts {
		if c, ok := stmt.(*sqlast.CreateTableStmt); ok {
			xcreateTable, err := ast.XCreateTableTo(c)
			if err != nil {
				return nil, err
			}
			createTable := light.CreateTableTo(xcreateTable)
			table := fromCreateTable(createTable)
			tables[table.TableName] = table
		} else if c, ok := stmt.(*sqlast.AlterTableStmt); ok {
			xalterTable, err := ast.XAlterTableTo(c)
			if err != nil {
				return nil, err
			}
			alterTable := light.AlterTableTo(xalterTable)
			tname := strings.Join(alterTable.TableName.Idents, ".")
			if x := alterTable.Action.GetAddConstraintItem(); x != nil {
				if y := x.Spec.GetReferenceItem(); y != nil {
					if fks[tname] == nil {
						fks[tname] = make([]*godbi.Fk, 0)
					}
					expr := y.KeyExpr
					fks[tname] = append(fks[tname], &godbi.Fk{
						FkTable:  expr.TableName,
						FkColumn: expr.Columns[0],
						Column:   y.Columns[0]})
				}
			}
		}
	}

	team := database{DBDriver: godbi.Postgres}
	postgres := &postgresIO{database: team, tables: tables, fks: fks}
	postgres.database.schema = postgres
	return postgres, nil
}

func dataTypeToTypeName(dt *xlight.Type) string {
	if x := dt.GetIntData(); x != nil {
		return "int"
	} else if x := dt.GetBigIntData(); x != nil {
		return "int"
	} else if x := dt.GetSmallIntData(); x != nil {
		return "int"
	} else if x := dt.GetDecimalData(); x != nil {
		return "real"
	}
	return "string"
}

func fromCreateTable(createTable *xlight.CreateTableStmt) *godbi.Table {
	var pk, idauto string
	var cols []*godbi.Col
	var uniques []string

	for _, item := range createTable.Elements {
		if x := item.GetColumnDefElement(); x != nil {
			col := &godbi.Col{
				ColumnName: x.Name,
				Label:      x.Name}
			if y := x.DataType.GetCustomData(); y != nil {
				if strings.Join(y.Idents, ".") == "SERIAL" {
					col.TypeName = "int"
					col.Auto = true
					col.Notnull = true
					idauto = x.Name
				}
			} else if y := x.DataType.GetTimestampData(); y != nil {
				col.TypeName = "string"
				col.Auto = true
				col.Notnull = true
			} else {
				col.TypeName = dataTypeToTypeName(x.DataType)
			}
			for _, constraint := range x.Constraints {
				spec := constraint.Spec
				if z := spec.GetNotNullItem(); z != xlight.NotNullColumnSpec_NotNullColumnSpecUnknown {
					col.Notnull = true
				} else if z := spec.GetUniqueItem(); z != nil {
					if z.IsPrimaryKey {
						pk = x.Name
					}
					uniques = append(uniques, x.Name)
				}
			}
			cols = append(cols, col)
		}
	}

	return &godbi.Table{
		TableName: strings.Join(createTable.Name.Idents, "."),
		Columns:   cols,
		Pks:       []string{pk},
		IDAuto:    idauto}
}

func (p *postgresIO) tableNames(_ *sql.DB) ([]string, error) {
	var names []string
	for k := range p.tables {
		names = append(names, k)
	}
	return names, nil
}

func (p *postgresIO) getTable(_ *sql.DB, tableName string) (*godbi.Table, error) {
	return p.tables[tableName], nil
}

func (p *postgresIO) getFks(_ *sql.DB, tableName string) ([]*godbi.Fk, error) {
	return p.fks[tableName], nil
}
