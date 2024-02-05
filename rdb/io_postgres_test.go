package rdb

import (
	"os"
	"testing"
)

func TestPostgresIO(t *testing.T) {
	file, err := os.Open("people_postgres.pb.sql")
	if err != nil {
		t.Fatal(err)
	}
	dbname := "tutorial"
	object, err := newPostgresIO(dbname, file)
	if err != nil {
		t.Fatal(err)
	}

	molecule, err := object.GetMolecule(nil)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(dbname+".json", []byte(molecule.String()), 0666)
	if err != nil {
		t.Fatal(err)
	}
	/*
		data, err := os.ReadFile(dbname + ".json")
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != molecule.String() {
			t.Errorf("not equal: %s", molecule.String())
		}
	*/
}
