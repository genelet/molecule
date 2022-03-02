package main

import (
	"context"
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

	db.Exec(`drop table if exists m_b`)
	db.Exec(`drop table if exists m_a`)
	db.Exec(`CREATE TABLE m_a (
		id int auto_increment not null primary key,
		x varchar(8), y varchar(8), z varchar(8))`)
	db.Exec(`CREATE TABLE m_b (
		tid int auto_increment not null primary key,
		child varchar(8),
		id int)`)

	ctx := context.Background()
	METHODS := map[string]string{"LIST": "topics", "GET": "edit", "POST": "insert", "PUT": "update", "PATCH": "insupd", "DELETE": "delete"}

	molecule, err := godbi.NewMoleculeJson([]byte(`{"atoms":
[{
	"tableName": "m_a",
	"pks": [ "id" ],
	"idAuto": "id",
    "columns": [
		{"columnName":"x", "label":"x", "typeName":"string", "notnull":true },
		{"columnName":"y", "label":"y", "typeName":"string", "notnull":true },
		{"columnName":"z", "label":"z", "typeName":"string" },
		{"columnName":"id", "label":"id", "typeName":"int", "auto":true }
    ],
    "uniques":["x","y"],
	"actions": [{
		"actionName": "insupd",
		"isDo": true,
		"nextpages": [{
			"tableName": "m_b",
			"actionName": "insert",
			"relateArgs": { "id": "id" },
			"marker": "m_b"
		}]
	},{
		"actionName": "insert",
		"isDo": true,
		"nextpages": [{
			"tableName": "m_b",
			"actionName": "insert",
			"relateArgs": { "id": "id" },
			"marker": "m_b"
		}]
	},{
		"actionName": "edit",
		"nextpages": [{
			"tableName": "m_b",
			"actionName": "topics",
			"relateExtra": { "id": "id" }
		}]
	},{
		"actionName": "delete",
		"prepares": [{
			"tableName": "m_b",
			"actionName": "delecs",
			"relateArgs": { "id": "id" }
		}]
	},{
		"actionName": "topics",
		"nextpages": [{
			"tableName": "m_a",
			"actionName": "edit",
			"relateExtra": { "id": "id" }
		}]
	}]
},{
	"tableName": "m_b",
	"pks": [ "tid" ],
	"fks": [{"fkTable":"m_a", "fkColumn":"id", "column":"id"}],
	"idAuto": "tid",
	"columns": [
		{"columnName":"tid", "label":"tid", "typeName":"int", "notnull": true, "auto":true},
		{"columnName":"child", "label":"child", "typeName":"string"},
		{"columnName":"id", "label":"id", "typeName":"int", "notnull": true}
	],
	"actions": [{
		"isDo": true,
		"actionName": "insert"
	},{
		"actionName": "edit"
	},{
		"isDo": true,
		"actionName": "delete"
	},{
		"isDo": true,
		"actionName": "delecs",
		"nextpages": [{
			"tableName": "m_b",
			"actionName": "delete",
			"relateArgs": { "tid": "tid" }
		}]
	},{
		"actionName": "topics"
	}]
}]
}`))

	var lists []interface{}

	// the 1st web requests creates id=1 to the m_a and m_b tables:
	//
	args := map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john", "m_b": []map[string]interface{}{{"child": "john"}, {"child": "john2"}}}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args)
	if err != nil { panic(err) }

	// the 2nd request just updates, becaues [x,y] is unique in m_a.
	// but creates a new record in tb for id=1
	//
	args = map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "m_b": map[string]interface{}{"child": "sam"}}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args)
	if err != nil { panic(err) }

	// the 3rd request creates id=2
	//
	args = map[string]interface{}{"x": "c1234567", "y": "d1234567", "z": "e1234", "m_b": map[string]interface{}{"child": "mary"}}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args)
	if err != nil { panic(err) }

	// the 4th request creates id=3
	//
	args = map[string]interface{}{"x": "e1234567", "y": "f1234567", "z": "e1234", "m_b": map[string]interface{}{"child": "marcus"}}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args)
	if err != nil { panic(err) }

    // GET all
    args = map[string]interface{}{}
    lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
    if err != nil { panic(err) }
	fmt.Printf("Step 1: %v\n", lists)

    // GET one
    args = map[string]interface{}{"id": 1}
    lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["GET"], args)
    if err != nil { panic(err) }
	fmt.Printf("Step 2: %v\n", lists)

    // DELETE
    lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["DELETE"], map[string]interface{}{"id": 1})
	if err != nil { panic(err) }

    // GET all m_a
    args = map[string]interface{}{}
    lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
    if err != nil { panic(err) }
	fmt.Printf("Step 3: %v\n", lists)

    // GET all m_b
    args = map[string]interface{}{}
    lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"], args)
    if err != nil { panic(err) }
	fmt.Printf("Step 4: %v\n", lists)

    db.Exec(`drop table if exists m_a`)
    db.Exec(`drop table if exists m_b`)
}
