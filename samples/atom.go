package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/genelet/molecule/godbi"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbName := os.Getenv("DBNAME")
    db, err := sql.Open("mysql", dbUser + ":" + dbPass + "@/" + dbName)
    if err != nil { panic(err) }
    defer db.Close()

	db.Exec(`CREATE TABLE m_a (
		id int auto_increment not null primary key,
        x varchar(8), y varchar(8), z varchar(8))`)

	str := `{
    "tableName":"m_a",
    "pks":["id"],
    "idAuto":"id",
    "columns": [
		{"columnName":"x", "label":"x", "typeName":"string", "notnull":true },
		{"columnName":"y", "label":"y", "typeName":"string", "notnull":true },
		{"columnName":"z", "label":"z", "typeName":"string" },
		{"columnName":"id", "label":"id", "typeName":"int", "auto":true }
    ],
	"uniques":["x","y"],
	"actions": [
		{ "isDo":true, "actionName": "insert" },
		{ "isDo":true, "actionName": "insupd" },
		{ "isDo":true, "actionName": "delete" },
		{ "actionName": "topics" },
		{ "actionName": "edit" }
	]}`
	atom, err := godbi.NewAtomJson([]byte(str))
	if err != nil { panic(err) }

	var lists []interface{}
	// the 1st web requests is assumed to create id=1 to the m_a table
	//
	args := map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "temp"}
	lists, err = atom.RunAtom(db, "insert", args)
	if err != nil { panic(err) }

	// the 2nd request just updates, becaues [x,y] is defined to the unique
	//
	args = map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "zzzzz"}
	lists, err = atom.RunAtom(db, "insupd", args)
	if err != nil { panic(err) }

	// the 3rd request creates id=2
	//
	args = map[string]interface{}{"x": "c1234567", "y": "d1234567", "z": "e1234"}
	lists, err = atom.RunAtom(db, "insert", args)
	if err != nil { panic(err) }

	// the 4th request creates id=3
	//
	args = map[string]interface{}{"x": "e1234567", "y": "f1234567", "z": "e1234"}
	lists, err = atom.RunAtom(db, "insupd", args)
	if err != nil { panic(err) }

	// GET all
	args = map[string]interface{}{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil { panic(err) }
	fmt.Printf("Step 1: %v\n", lists)

	// GET one
	args = map[string]interface{}{"id": 1}
	lists, err = atom.RunAtom(db, "edit", args)
	if err != nil { panic(err) }
	fmt.Printf("Step 2: %v\n", lists)

	// DELETE
	args = map[string]interface{}{"id": 1}
	lists, err = atom.RunAtom(db, "delete", args)
	if err != nil { panic(err) }

	// GET all
	args = map[string]interface{}{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil { panic(err) }
	fmt.Printf("Step 3: %v\n", lists)

	db.Exec(`drop table if exists m_a`)
}
