# molecule

_molecule_ runs complex RESTful actions on selective data fields in related database tables like gRPC/GraphQL. The relationships between tables, which are usually described in JSON, include logic operators, data filters, action triggers and SQL foreign-key constraints etc. 

In _molecule_, a table and its associated actions build up an _atom_. Atoms and relationships between them build up a _molecule_.

While traditional REST acts on individual table, _molecule_ acts on whole RDB across all tables.

This package has pre-defined 6 RESTful actions, with which we can run most database tasks with little or no coding. For example, think about a gRPC application. We can create a Postgres database representing data stream's protocol buffer, and a JSON config representing relationships between the tables (which are usually mapped to _protobuf messages_). With _molecule_, we can process gRPC's input and output calls at once. Beneath the surface, _molecule_ will handle detailed reads and writes on proper tables with given logic. 

Check *godoc* for package details:

[![GoDoc](https://godoc.org/github.com/genelet/molecule?status.svg)](https://godoc.org/github.com/genelet/molecule)

The package is fully tested for PostgreSQL, MySQL and SQLite.

<br /><br />

## Chapter 1. INSTALLATION

To install:

> $ go get -u github.com/genelet/molecule

To use _molecule_, in your GO program:

```go
import ("github.com/genelet/molecule/godbi")
```

You may also encode and process _molecule_ in protocol buffer, which we have named [Graph](https://godoc.org/github.com/genelet/molecule/gometa). To use _Graph_:

```go
import ("github.com/genelet/molecule/gometa")
```

#### 1.1) Termilogy

- *table*: a database table;
- *action*: a *SELECT* or *DO* database action;
- *atom*: a table with actions (or *node* in *graph*); 
- *connection*: a relationship between atoms (or *edge*);
- *molecule*: whole set of atoms which act with each other in relationship (or *graph*);
- *RDB*: relationational database system.

#### 1.2) Arguments

The following names in functions are defined to be:

Name | Type | IN/OUT | Meaning
---- | ---- | ------ | -------
*args* | `...interface{}` | IN | arguments
*ARGS* | `map[string]interface{}` | IN | input data
*extra* | `...map[string]interface{}` | IN | _WHERE_ constraints
*lists* | `[]map[string]interface{}` | OUT | output data

#### 1.3) Three Levels of Usages:

- Basic Usage: run raw SQL statements using DBI struct;
- Atom Usage: run actions on single table; 6 RESTful actions are pre-defined in this package;
- Molecule Usage: run GraphQL/gRPC actions of multiple relational atoms.


<br /><br />

## Chapter 2. BASIC USAGE

In this example, we create table _letters_ with 3 rows, then search and put the data into *lists*.

<details>
    <summary>Click for Basic Usage Sample</summary>
    <p>

```go
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
```

Running it will output

```bash
[map[id:1 x:m] map[id:2 x:n] map[id:3 x:p]]
```

</p>
</details>

Here are frequently used functions in the `DBI usage`. For detailed list, please check the document.

### 2.1) DBI

The `DBI` type embeds the standard SQL handle. The extra *LastID* field is used to stored auto-generated series number, if the database supports it.

```go
package molecule

type DBI struct {
    *sql.DB          
    LastID   int64  // saves the last inserted id, if any
}

```

### 2.2) DoSQL

The same as DB's `Exec`, except it returns error.

```go
func (*DBI) DoSQL(query string, args ...interface{}) error
```

### 2.3) TxSQL

The same as _DoSQL_ but using transaction.

```go
func (*DBI) TxSQL(query string, args ...interface{}) error
```

### 2.4) Select

Return query data into *lists*, with data types determined dynamically.

```go
Select(lists *[]interface{}, query string, args ...interface{}) error
```

### 2.5) SelectSQL

The same as *Select* but using pre-defined labels and types.

```go
func (*DBI) SelectSQL(lists *[]map[string]interface{}, labels []interface{}, query string, args ...interface{}) error
```

<details>
    <summary>Click for Example</summary>

The following example assigns key names _TS_, _id_, _Name_, _Length_, _Flag_ and _fv_, of data types _string_, _int_, _string_, _int8_, _bool_ and _float32_, to the returned rows:

```go
lists := make([]map[string]interface{})
err = molecule.SelectSQL(
    &lists, 
    `SELECT ts, id, name, len, flag, fv FROM mytable WHERE id=?`,
    []interface{}{[2]string{"TS","string"], [2]string{"id","int"], [2]string{"Name","string"], [2]string{"Length","int8"], [2]string{"Flag","bool"], [2]string{"fv","float32"]},
    1234)
```
	    
It outputs:
	    
```json
    {"TS":"2019-12-15 01:01:01", "id":1234, "Name":"company", "Length":30, "Flag":true, "fv":789.123},
```

</details>

### 2.6) GetSQL

If there is only one row returned, use this function to get a map.

```go
func (*DBI) GetSQL(res map[string]interface{}, query string, labels []interface{}, args ...interface{}) error
```

<br /><br />


## Chapter 3. ATOM USAGE

In the following example, we define table, columns and actions in JSON, and run REST actions on the table. 

<details>
	<summary>Click for Atom Usage Sample: Default and Customized Actions</summary>
	
```go
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
```

Running the program will result in

```bash
Step 1: [map[id:1 x:a y:b z:c]]
Step 2: [map[id:2 x:c y:d z:c]]
Step 3: [map[id:1 x:a y:b z:c] map[id:2 x:c y:d z:c]]
Step 4: [map[x:c y:z]]
Step 5: [map[id:2 x:c y:z z:c]]
Step 6: [map[z:c]]
```
	
</details>

<details>
	<summary>Click for Atom Usage Sample: Atom</summary>
	
```go
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

```

Running the program will result in

```bash
Step 1: [map[id:1 x:a1234567 y:b1234567 z:zzzzz] map[id:2 x:c1234567 y:d1234567 z:e1234] map[id:3 x:e1234567 y:f1234567 z:e1234]]
Step 2: [map[id:1 x:a1234567 y:b1234567 z:zzzzz]]
Step 3: [map[id:2 x:c1234567 y:d1234567 z:e1234] map[id:3 x:e1234567 y:f1234567 z:e1234]]
```
	
</details>

Here are frequently used data types and functions in the atom usage.

<br />

### 3.1) Col

Define a table column type:

```go
type Col struct {
    ColumnName string  `json:"columnName" hcl:"columnName"`
    TypeName string    `json:"typeName" hcl:"typeName"`
    Label string       `json:"label" hcl:"label"`
    Notnull bool       `json:"notnull" hcl:"notnull"`
    Auto bool          `json:"auto" hcl:"auto"`
    // true for a one-to-may recurse column
    Recurse bool       `json:"recurse,omitempty" hcl:"recurse,optional"`
}
```
where _ColumnName_ is the column name. _TypeName_ is column's type. _Lable_ is the label for the column. _Notnull_ marks if the column can't be null. _Auto_ means if the column can be automatically assigned with value e.g. timestamp, auto id etc. And _Recurse_ means it recursively references to table's primary key in one-to-many relation.

### 3.2) Fk

SQL's foreign key. It's a relationship between 2 atoms:

```go
type Fk struct {
    FkTable  string    `json:"fkTable" hcl:"fkTable"`
    FkColumn string    `json:"fkColumn" hcl:"fkColumn"`
    Column   string    `json:"column" hcl:"column"`
}
```

where _FkTable_ means a foreign table, _FkColumn_ foreign table's column, and _Column_ the column in the current table. 

- Foreign key is defined to be a single column in this package.
- Foreign key can be defined even if there is no native SQL foreign key, like NoSQL or time-series database.
- Foreign key is only used in action *Delecs*.

### 3.3) Table

_Table_ describes a database table.

```go
type Table struct {
    TableName string   `json:"tableName" hcl:"tableName"`
    Columns   []*Col   `json:"columns" hcl:"columns"`
    Pks       []string `json:"pks,omitempty" hcl:"pks,optional"`
    IdAuto    string   `json:"idAuto,omitempty" hcl:"idAuto,optional"`
    Fks       []*Fk    `json:"fks,omitempty" hcl:"fks,optional"`
    Uniques   []string `json:"uniques,omitempty" hcl:"uniques,optional"`
}

```

where _TableName_ is the table name. _Columns_ are all columns. _Pks_ is the primary key. _IdAuto_ is the auto ID. _Fks_ is a list of foreign-key relationships. And _Uniques_ is the combination of columns uniquely defining the row.

### 3.4) Connection

_connection_ is associated with specific table. It defines relationship to another table for how data are passed under a specific action.

```go
type Connection struct {
    TableName   string            `json:"tableName" hcl:"tableName,label"`
    ActionName  string            `json:"actionName" hcl:"actionName,label"`
    RelateArgs  map[string]string `json:"relateArgs,omitempty" hcl:"relateArgs"`
    RelateExtra map[string]string `json:"relateExtra,omitempty" hcl:"relateExtra"`
    Dimension  ConnectType        `json:"dimension,omitempty" hcl:"dimension,label"`
    Marker     string             `json:"marker,omitempty" hcl:"marker,label"`
}
```

where _TableName_ is the database table name. _ActionName_ is the action name. _RelateArgs_ is a filter, that maps an output data from this table to the input data of the next table. _RelateExtra_ is for _where_ constraint. _Dimension_ is a relation type. And _Marker_ is a string marker. For an input action like _insert_ and _insupd_, _Marker_ will be used as key to reference row data; and for an output action like _topics_, it stores the row under the marker.

### 3.5) Action

*Action* defines an action, such as *CRUD*, on a table. It implements interface _Capability_.

```go
type Action struct {
    ActionName string `json:"actionName,omitempty" hcl:"actionName,optional"`
    Prepares  []*Connection `json:"Prepares,omitempty" hcl:"Prepares,block"`
    Nextpages []*Connection `json:"nextpages,omitempty" hcl:"nextpages,block"`
    IsDo      bool          `json:"isDo,omitempty" hcl:"isDo,optional"`
}
```

where _Prepares_ is a list of actions to run before the current action, and _Nextpages_ actions to follow. We can think _Prepares_ as `pre triggers` and _Nextpages_ as `post triggers` in standard SQL. _IsDo_ indicates if it is associated with SQL's _DO_ query.

### 3.6) RunActionContext

This is the main function in _Capability_. It takes input data _ARGS_ and optional constraint _extra_, and run. The output is a slice of interface, and an optional error.

```go
RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error)
```

### 3.7) Atom

An atom is made of a table and its pre-defined actions. 

```go
type Atom struct {
 Table
 Actions []Capability `json:"actions,omitempty" hcl:"actions,optional"`
}
```

### 3.8) RunAtomContext

This is the main function for _action_. It takes input data _ARGS_ and optional constraint _extra_, and acts. The output is a slice of interface, and an optional error.
```go
RunActomContext(ctx context.Context, db *sql.DB, action string, ARGS interface{}, extra ...map[string]interface{}) ([]interface{}, error)
```

<br /><br />

## Chapter 4. RESTFUL ACTIONS

_molecule_ has defined the following 7 RESTful actions. Most RDB tasks can be fulfilled using these actions.


### 4.1) Insert

Add a new row into table.

```go
type Insert struct {
    Action
}
```

### 4.2) Update

Update a row by primary key.

```go
type Update struct {
    Action
    Empties  []string  `json:"empties,omitempty" hcl:"empties,optional"`
}
```

where _Empties_ defines columns
which will be forced to be empty or null if having no input data.

### 4.3) Insupd

If insert a new row if the row is unique, otherwise update.

```go
type Insupd struct {
    Action
}
```

### 4.4) Edit

Query one row by primary key. We can dynamically pass the field names as concated value in _ARGS[FIELDS]_. _Joints_ is provided optionally so a single _Edit_ could run a more sophisticated _JOIN_ statement.

```go
type Edit struct {
    Action
    Joins    []*Joint   `json:"joins,omitempty" hcl:"join,block"`
    FIELDS   string    `json:"fields,omitempty" hcl:"fields"`
}
```

### 4.5) Topics

```go
type Topics struct {
    Action
    Joints []*Joint    `json:"joints,omitempty" hcl:"joints,block"`
    FIELDS string      `json:"fields,omitempty" hcl:"fields"`

    TotalForce  int    `json:"total_force,omitempty" hcl:"total_force,optional"`
    MAXPAGENO   string `json:"maxpageno,omitempty" hcl:"maxpageno,optional"`
    TOTALNO     string `json:"totalno,omitempty" hcl:"totalno,optional"`
    ROWCOUNT    string `json:"rawcount,omitempty" hcl:"rawcount,optional"`
    PAGENO      string `json:"pageno,omitempty" hcl:"pageno,optional"`
    SORTBY      string `json:"sortby,omitempty" hcl:"sortby,optional"`
    SORTREVERSE string `json:"sortreverse,omitempty" hcl:"sortreverse,optional"`
}
```

where pagination is defined by:

Field | Default | Meaning in Input Data `ARGS`
--------- | ------- | -----------------------
_MAXPAGENO_ | "maxpageno" | how many pages in total
_TOTALNO_ | "totalno" | how many records in total
_ROWCOUNT_ | "rowcount" | how many record in each page
_PAGENO_ | "pageno" | return only data of the specific page
_SORTBY_ | "sortby" | sort the returned data by this
_SORTREVERSE_ | "sortreverse" | 1 to return the data in reverse

and _TotalForce_ is: 0 for not calculating total number of records; -1 for calculating; and 1 for optionally calculating. In the last case, if there is no input data for `ROWCOUNT` or `PAGENO`, there is no pagination information.

### 4.6) Delete

Delete a row by primary key.

```go
type Delete struct {
    Action
}
```

### 4.7) Delecs

Get all keys and foreign keys for a row. In _molecule_, when to delete a row in this table, deletions could be triggered in other tables. We always run _Delecs_ before _Delete_ so as to get keys ready for the related tables.

```go
type Delecs struct {
    Action
}
```

<br /><br />

## Chapter 5. MOLECULE USAGE


*Molecule* is a collection of all atoms in the database and implements the *Run* function to react.

```go
type Molecule struct {
    Atoms []Navigate `json:"atoms" hcl:"atoms"`
    DatabaseName string `json:"databaseName" hcl:"databaseName"`
    DBDriver DBType `json:"dbDriver" hcl:"dbDriver"`
	Stopper
}

type Stopper interface {
	Sign(tableObj *Table, item interface{}) bool
}
```

where _DBDriver_ is one of database drive defined:

```go
    SQLDefault DBType = iota
    SQLRaw
    SQLite
    MySQL
    Postgres
    TSMillisecond
    TSMicrosecond
```

Stopper stops molecule's chain actions at an earlier stage defined by _Sign_ is true.

<details>
    <summary>Click for Molecule Usage Sample</summary>
    <p>

```go
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
```

Running it will result in:

```bash
Step 1: [map[id:1 m_a_edit:[map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:john2 id:1 tid:2] map[child:sam id:1 tid:3]] x:a1234567 y:b1234567 z:zzzzz]] x:a1234567 y:b1234567 z:zzzzz] map[id:2 m_a_edit:[map[id:2 m_b_topics:[map[child:mary id:2 tid:4]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_a_edit:[map[id:3 m_b_topics:[map[child:marcus id:3 tid:5]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]
Step 2: [map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:john2 id:1 tid:2] map[child:sam id:1 tid:3]] x:a1234567 y:b1234567 z:zzzzz]]
Step 3: [map[id:2 m_a_edit:[map[id:2 m_b_topics:[map[child:mary id:2 tid:4]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_a_edit:[map[id:3 m_b_topics:[map[child:marcus id:3 tid:5]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]
Step 4: [map[child:mary id:2 tid:4] map[child:marcus id:3 tid:5]]
```

</p>
</details>

<br />

### 5.1) Construct

Use JSON to build a molecule

```go
func NewMoleculeJsonFile(filename string, cmap ...map[string][]Capability) (*Molecule, error) 
```

where _cmap_ is for customized actions not in the default list.

### 5.2) Run action on atom

We can run any action on any atom by names using _RunConext_. The output is data as a slice of interface, and an optional error.

```go
func (self *Molecule) RunContext(ctx context.Context, atom, action string, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]map[string]interface{}, error)
```

Unlike traditional REST, which is limited to a sinlge table and sinle action, _RunContext_ will act on related tables and trigger associated actions.

Please check [the full document](https://godoc.org/github.com/genelet/molecule) for usage.

