package rdb

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLite(t *testing.T) {
	var databaseName = "imdb"
	db, err := sql.Open("sqlite3", databaseName+".sqlite")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	object := newSQLite(databaseName)
	molecule, err := object.GetMolecule(db)
	if err != nil {
		t.Fatal(err)
	}

	/*
		err = os.WriteFile(databaseName + ".json", []byte(molecule.String()), 0666)
		if err != nil {
			t.Fatal(err)
		}
	*/

	data, err := os.ReadFile(databaseName + ".json")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != String(molecule) {
		t.Errorf("not equal: %s", String(molecule))
	}
}
