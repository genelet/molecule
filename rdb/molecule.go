package rdb

import (
	"database/sql"
	"fmt"

	"github.com/genelet/molecule/godbi"
)

// NewMolecule returns a molecule of database using primary keys, foreign keys
// and auto increment numbers etc. The arguments are:
//
// db: the standard database handler
// driver: the DBType defined in molecule. Currently only the three databases
// are supported, godbi.Postgres, godbi.MySQL and godbi.SQLite
// dbName: the name of database
// 
func NewMolecule(db *sql.DB, driver godbi.DBType, dbName string) (*godbi.Molecule, error) {
	switch driver {
	case godbi.MySQL:
		return newMySQL(dbName).GetMolecule(db)
	case godbi.SQLite:
		return newSQLite(dbName).GetMolecule(db)
	case godbi.Postgres:
		return newPostgres(dbName).GetMolecule(db)
	default:
	}
	return nil, fmt.Errorf("db type %v not found", driver)
}
