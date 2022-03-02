package main

import (
	"context"
    "database/sql"
    "fmt"
    "os"

    "github.com/genelet/molecule/godbi"
    _ "github.com/go-sql-driver/mysql"
)

type SQL struct {
    godbi.Action
    Statement string   `json:"statement"`
}

func (self *SQL) RunAction(db *sql.DB, t *godbi.Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
    return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *SQL) RunActionContext(ctx context.Context, db *sql.DB, t *godbi.Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
    lists := make([]interface{}, 0)
    dbi := &godbi.DBI{DB: db}
    err := dbi.SelectContext(ctx, &lists, self.Statement, ARGS["marker"])
    return lists, err
}

func main() {
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbName := os.Getenv("DBNAME")
    db, err := sql.Open("mysql", dbUser + ":" + dbPass + "@/" + dbName)
    if err != nil { panic(err) }
    defer db.Close()

    db.Exec(`CREATE TABLE testing (
		id int auto_increment,
		x varchar(255),
		y varchar(255),
		z varchar(255),
		primary key (id))`)

    table := &godbi.Table{
		TableName: "testing",
		Pks:[]string{"id"},
		IdAuto:"id",
		Columns: []*godbi.Col{
			&godbi.Col{ColumnName:"x",  Label:"x",  TypeName:"string", Notnull:true},
			&godbi.Col{ColumnName:"y",  Label:"y",  TypeName:"string", Notnull:true},
			&godbi.Col{ColumnName:"z",  Label:"z",  TypeName:"string"},
			&godbi.Col{ColumnName:"id", Label:"id", TypeName:"int", Auto:true},
		},
	}

    insert := &godbi.Insert{Action:godbi.Action{ActionName: "insert"}}
    topics := &godbi.Topics{Action:godbi.Action{ActionName: "topics"}}
    update := &godbi.Update{Action:godbi.Action{ActionName: "update"}}
    edit   := &godbi.Edit{Action:godbi.Action{ActionName: "edit"}}
    // a custom action
    sql   := &SQL{Action:godbi.Action{ActionName: "sql"}, Statement:"SELECT z FROM testing WHERE id=?"}

    args := map[string]interface{}{"x":"a","y":"b","z":"c"}
    lists, err := insert.RunAction(db, table, args)
    if err != nil { panic(err) }
    fmt.Printf("Step 1: %v\n", lists)

    args = map[string]interface{}{"x":"c","y":"d","z":"c"}
    lists, err = insert.RunAction(db, table, args)
    if err != nil { panic(err) }
    fmt.Printf("Step 2: %v\n", lists)

    lists, err = topics.RunAction(db, table, nil)
    fmt.Printf("Step 3: %v\n", lists)

    args = map[string]interface{}{"id":2,"x":"c","y":"z"}
    lists, err = update.RunAction(db, table, args)
    if err != nil { panic(err) }
    fmt.Printf("Step 4: %v\n", lists)

    args = map[string]interface{}{"id":2}
    lists, err = edit.RunAction(db, table, args)
    fmt.Printf("Step 5: %v\n", lists)

    args = map[string]interface{}{"marker":1}
    lists, err = sql.RunAction(db, table, args)
    fmt.Printf("Step 6: %v\n", lists)

    db.Exec(`DROP TABLE IF EXISTS testing`)
    os.Exit(0)
}
