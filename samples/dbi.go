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

    dbi := &godbi.DBI{DB:db}

    // create a new table and insert 3 rows
    //
    _, err = db.Exec(`CREATE TABLE letters (
        id int auto_increment primary key,
        x varchar(1))`)
    if err != nil { panic(err) }
    _, err = db.Exec(`insert into letters (x) values ('m')`)
    if err != nil { panic(err) }
    _, err = db.Exec(`insert into letters (x) values ('n')`)
    if err != nil { panic(err) }
    _, err = db.Exec(`insert into letters (x) values ('p')`)
    if err != nil { panic(err) }

    // select all data from the table using Select
    //
    lists := make([]interface{}, 0)
    err = dbi.Select(&lists, "SELECT id, x FROM letters")
    if err != nil { panic(err) }

    fmt.Printf("%v\n", lists)

    dbi.Exec(`DROP TABLE IF EXISTS letters`)
    db.Close()
}
