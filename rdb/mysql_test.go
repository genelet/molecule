package rdb

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQL(t *testing.T) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	databaseName := "classicmodels"
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+databaseName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	object := newMySQL(databaseName)
	molecule, err := object.GetMolecule(db)
	if err != nil {
		t.Fatal(err)
	}

	/*
		err = os.WriteFile(databaseName + "1.json", []byte(molecule.String()), 0666)
		if err != nil {
			t.Fatal(err)
		}
	*/

	data, err := os.ReadFile(databaseName + ".json")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != molecule.String() {
		t.Errorf("not equal: %s", molecule.String())
	}
}
