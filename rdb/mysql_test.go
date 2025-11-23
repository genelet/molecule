package rdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/genelet/horizon/dethcl"
	"github.com/genelet/molecule/godbi"
	_ "github.com/go-sql-driver/mysql"
)

func String(m *godbi.Molecule) string {
	bs, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", bs)
}

func HCLString(self *godbi.Molecule) string {
	bs, err := dethcl.Marshal(self)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", bs)
}

func TestMySQL(t *testing.T) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	databaseName := "classicmodels"
	dbHOST := os.Getenv("DBHOST")
	var db *sql.DB
	var err error
	if dbHOST != "" {
		db, err = sql.Open("mysql", dbUser+":"+dbPass+"@("+dbHOST+")/"+databaseName)
	} else {
		db, err = sql.Open("mysql", dbUser+":"+dbPass+"@/"+databaseName)
	}
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

	if string(data) != String(molecule) {
		t.Errorf("not equal: %s", String(molecule))
	}
}
