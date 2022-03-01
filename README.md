# molecule

_molecule_ runs RESTful actions on related database tables and selective data fields in RDB like gRPC/GraphQL. The relationship, which is usually described in JSON, include SQL constrains like foreign-key, and flexible logic operations like data filters and action triggers etc. In _molecule_, a table and its associated actions build up an _atom_. Atoms and relationships between atoms build up a _molecule_.

While traditional REST acts on individual table, _molecule_ acts on a whole database across all tables in the database.

This package has pre-defined 6 RESTful actions, with which we can run most database tasks with little or no coding. For example,  think about a gRPC application. With _molecule_, we can create a Postgres database representing data stream's protocol buffer, and a JSON config to represent relationships between Postgres tables (which are usually mapped to _protobuf messages_). Then _molecule_ will process input and output streams as a gRPC call at once. Beneath the surface, _molecule_ handles detailed reads and writes on proper tables with given logic. 

Check *godoc* for package details:

[![GoDoc](https://godoc.org/github.com/genelet/molecule?status.svg)](https://godoc.org/github.com/genelet/molecule)

_molecule_ has 3 levels of usages:

- _Basic_: run raw SQL statements;
- _Atom_: run actions on single table; 6 RESTful actions are pre-defined in this package;
- _Molecule_: run GraphQL/gRPC actions of multiple relational atoms.


The package is fully tested for PostgreSQL, MySQL and SQLite.


<br /><br />

## Chapter 1. INSTALLATION

To install:

> $ go get -u github.com/genelet/molecule

#### 1.1) Termilogy

- *table*: a database table;
- *action*: a *SELECT* or *DO* database action;
- *atom*: a table with actions (or a *node* in *graph*); 
- *connection*: a relationship between atoms (or *edge*);
- *molecule*: whole set of atoms which act with each other in relationship (or *graph);
- *RDB*: relationational database system (or *meta*)

#### 1.2) Arguments

The following names in functions are defined to be:

Name | Type | IN/OUT | Meaning
---- | ---- | ------ | -------
*args* | `...interface{}` | IN | arguments
*ARGS* | `map[string]interface{}` | IN | input data
*extra* | `...map[string]interface{}` | IN | _WHERE_ constraints
*lists* | `[]map[string]interface{}` | OUT | output data

<br /><br />

## Chapter 2. BASIC USAGE

In this example, we create table _letters_ with 3 rows, then search and put the data into *lists*.

<details>
    <summary>Click for DBI example</summary>
    <p>

```go
package main

import (
    "log"
    "database/sql"
    "os"
    "github.com/genelet/molecule"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbName := os.Getenv("DBNAME")
    db, err := sql.Open("mysql", dbUser + ":" + dbPass + "@/" + dbName)
    if err != nil { panic(err) }
    defer db.Close()

    dbi := &molecule.DBI{DB:db}

    // create a new table and insert 3 rows
    //
    _, err = db.Exec(`DROP TABLE IF EXISTS letters`)
    if err != nil { panic(err) }
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
    lists := make([]map[string]interface{}, 0)
    err = dbi.Select(&lists, "SELECT id, x FROM letters")
    if err != nil { panic(err) }

    log.Printf("%v", lists)

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

Here are basic data type and functions in the DBI usage. For details, please check the [document](https://)

### 2.1  _DBI_

The `DBI` type embeds the standard SQL handle.

```go
package molecule

type DBI struct {
    *sql.DB          
    LastID   int64  // saves the last inserted id, if any
}

```

To create a new handle

```go
dbi := &DBI{DB: the_standard_sql_handle}
```

<br />

### 2.2  _DoSQL_

```go
func (*DBI) DoSQL(query string, args ...interface{}) error
```

It is the same as DB's `Exec`, except it returns error only.

<br />

### 2.3  _TxSQL_

```go
func (*DBI) TxSQL(query string, args ...interface{}) error
```

The same as _DoSQL_ using transaction.

<br />

### 2.4   _Select_

```go
func (*DBI) Select(lists *[]map[string]interface{}, query string, args ...interface{}) error
```

It runs a *Select* query and saves the result into *lists*. The data types of the rows are determined dynamically by the generic SQL handle.

<details>
    <summary>Click for example</summary>
    <p>

```go
lists := make([]map[string]interface{})
err = dbi.Select(&lists,
    `SELECT ts, id, name, len, flag, fv FROM mytable WHERE id=?`, 1234)
```

will select all rows with *id=1234* into *lists*:

```json
    {"ts":"2019-12-15 01:01:01", "id":1234, "name":"company", "len":30, "flag":true, "fv":789.123},
    ....
```

</p>
</details>

### 2.5) `SelectSQL`

```go
func (*DBI) SelectSQL(lists *[]map[string]interface{}, labels []interface{}, query string, args ...interface{}) error
```

The same as *Select*  but using *labels*, which can be  *string*
which defines renamed key, or *[2]string* which defines renamed key and its data type such as _int_, _int64_ and _string_ etc.

<details>
    <summary>Click for example</summary>
    <p>

The following example assigns key names _TS_, _id_, _Name_, _Length_, _Flag_ and _fv_, of data types _string_, _int_, _string_, _int8_, _bool_ and _float32_, to the returned rows:

```go
lists := make([]map[string]interface{})
err = dbi.querySQLLabel(&lists, 
    `SELECT ts, id, name, len, flag, fv FROM mytable WHERE id=?`,
 []interface{}{[2]string{"TS","string"], [2]string{"id","int"], [2]string{"Name","string"], [2]string{"Length","int8"], [2]string{"Flag","bool"], [2]string{"fv","float32"]},
    1234)
```

```json
    {"TS":"2019-12-15 01:01:01", "id":1234, "Name":"company", "Length":30, "Flag":true, "fv":789.123},
```

</p>
</details>

<br />

### 2.6)  _GetSQL_

If there is only one row returned, use this function to get a map.

```go
func (*DBI) GetSQL(res map[string]interface{}, query string, labels []interface{}, args ...interface{}) error
```

<br /><br />


## Chapter 3. ATOM USAGE

In this example, we define table, columns and actions in JSON, and let _molecule_ to run the REST actions. 

<details>
	<summary>Click here to see how atom works</summary>
	
```go
package main

import (
    "log"
    "os"
    "database/sql"
    "github.com/genelet/molecule"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbName := os.Getenv("DBNAME")
    db, err := sql.Open("mysql", dbUser + ":" + dbPass + "@/" + dbName)
    if err != nil { panic(err) }
    defer db.Close()

    db.Exec(`DROP TABLE IF EXISTS testing`)
    db.Exec(`CREATE TABLE testing (id int auto_increment, x varchar(255), y varchar(255), primary key (id))`)

    table := &molecule.Table{TableName: "testing", Pks:[]string{"id"}, IdAuto:"id"}

    insert := &molecule.Insert{Columns: []string{"x","y"}}
    topics := &molecule.Topics{Columns: map[string][]string{"id":{"id","int"}, "x":{"x","string"},"y":{"y","string"}}}
    update := &molecule.Update{Columns: []string{"id","x","y"}}
    edit   := &molecule.Edit{Columns: map[string][]string{"id":{"id","int"}, "x":{"x","string"},"y":{"y","string"}}}

    args := map[string]interface{}{"x":"a","y":"b"}
    lists, _, err := insert.RunAction(db, table, args)
    if err != nil { panic(err) }
    log.Println(lists)

    args = map[string]interface{}{"x":"c","y":"d"}
    lists, _, err = insert.RunAction(db, table, args)
    if err != nil { panic(err) }
    log.Println(lists)

    args = map[string]interface{}{}
    lists, _, err = topics.RunAction(db, table, args)
    log.Println(lists)

    args = map[string]interface{}{"id":2,"x":"c","y":"z"}
    lists, _, err = update.RunAction(db, table, args)
    if err != nil { panic(err) }
    log.Println(lists)

    args = map[string]interface{}{"id":2}
    lists, _, err = edit.RunAction(db, table, args)
    log.Println(lists)

    os.Exit(0)
}
```

Running the program will result in

```bash
[map[id:1 x:a y:b]]
[map[id:2 x:c y:d]]
[map[id:1 x:a y:b] map[id:2 x:c y:d]]
[map[id:2 x:c y:z]]
[map[id:2 x:c y:z]]
```

</p>
</details>

Here are basic data types and functions in the atom usage.

<br />

### 3.1) Col

Define a GO column type:

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
where _ColumnName_ is the name of column, _TypeName_ is the name of column type, _Lable_ is the label for the column, _Notnull_ marks if the column can't be null, _Auto_ means if the column can be automatically assigned with a value like timestamp, auto increment etc., and _Recurse_ means a special column it recursively reference to table's primary key in a one-to-many relation.

### 3.2) Fk

*foreign key Fk* is a relation between atoms:

```go
type Fk struct {
    FkTable  string    `json:"fkTable" hcl:"fkTable"`
    FkColumn string    `json:"fkColumn" hcl:"fkColumn"`
    Column   string    `json:"column" hcl:"column"`
}
```

where _FkTable_ means a foreign table, _FkColumn_ foriegn table's column and _Column_ the column in the current table. _Fk_ is similar to SQL's standard foreign key but 1) the columns are limited to be single, and 2) even if there is no SQL foreign key, like noSQL database or time-series database, we can still define them in atom.

### 3.3)  Table

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

where _TableName_ is the table name; _Columns_ are all columns, _Pks_ the primary key; _IdAuto_ the column having a auto increment series number, _Fks_ the foreign key relationship to other atoms, and _Uniques_ the combination of columns uniquely defining this row.

<br />

### 3.4) Connection

_connection_ defines data are passed between atoms.

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

where _TableName_ is the database table name. _ActionName_ is the action name. _RelateArgs_ maps an output data of this table to the input data of the next table. _RelateExtra_ is for _where_ constraint. _Dimension_ is a relation type. And _Marker_ is a string marker. During an input action such as _insert_ and _insupd_, _Marker_ points to this table's data; and during an output action such as _topics_, it place the table data under this string.

### 3.5) Action

*Action* defines an action, such as *CRUD*, on table. It should implement *Capability* interface:

```go
type Action struct {
    ActionName string `json:"actionName,omitempty" hcl:"actionName,optional"`
    Prepares  []*Connection `json:"Prepares,omitempty" hcl:"Prepares,block"`
    Nextpages []*Connection `json:"nextpages,omitempty" hcl:"nextpages,block"`
    IsDo      bool          `json:"isDo,omitempty" hcl:"isDo,optional"`
    Appendix  interface{}   `json:"appendix,omitempty" hcl:"appendix,block"`
}
```

where _Prepares_ is a list of actions to run before the current action, and _Nextpages_ actions to follow. 

_molecule_ has defined the following 6 RESTful actions as defaults for all tables in database.

#### 3.5.1) *Insert*

```go
type Insert struct {
    Action
    Columns    []string      `json:"columns,omitempty" hcl:"columns,optional"`
}
```

It adds a new row into _Columns_ of table. The input data are in *ARGS* defined in _RunActionContext_.

#### 2.2.2) *Update*

```go
ype Update struct {
    Action
    Columns    []string      `json:"columns,omitempty" hcl:"columns,optional"`
    Empties    []string      `json:"empties,omitempty" hcl:"empties,optional"`
}
```

It updates row's *Columns* using the primary key. _Empties_ defines columns
which are forced to be empty or null if having no input data.

#### 2.2.3) *Insupd*

```go
type Insupd struct {
    Action
    Columns    []string      `json:"columns,omitempty" hcl:"columns,optional"`
    Uniques    []string      `json:"uniques,omitempty" hcl:"uniques,optional"`
}
```

It checks if the input data is unique using input data from columns _Uniques_. If it exists, run a *Update* otherwise *Insert*.

#### 2.2.4) *Edit*

```go
type Edit struct {
    Action
    Joins    []*Join             `json:"joins,omitempty" hcl:"join,block"`
    Columns   map[string][]string `json:"columns" hcl:"columns"`
    FIELDS   string              `json:"fields,omitempty" hcl:"fields"`
}
```

```go
type Join struct {
	Name   string `json:"name" hcl:"name,label"`
	Alias  string `json:"alias,omitempty" hcl:"alias,optional"`
	Type   string `json:"type,omitempty" hcl:"type,optional"`
	Using  string `json:"using,omitempty" hcl:"using,optional"`
	On     string `json:"on,omitempty" hcl:"on,optional"`
	Sortby string `json:"sortby,omitempty" hcl:"sortby,optional"`
}

It searches one row from a table or joint tables, using main table's primary key.

*Joins* defines the joint tables. The main table is the first element
in the list. *Join* defines the table and how it is joint. *Name* is main table's name; *Alias* the table name alias; *Type* the join type such as *JOIN*, "INNER JOIN", "LEFT JOIN" etc.; "Using" the *USING* statement in SQL; "On" the "ON* statement; and *Sortby* which appears only in the main table the column name to sort the data.
 
In *Columns*, the keys are SQL column names and keys' values are renamed labels. See 1.4.2).

*FIELDS* defines a key in the input data which, if exists, tells us
to return only a few selected columns instead of whole set of columns.
By default, _FIELDS ="fields"_.  For example, in order to return
just *user_id* and *username*, we can assign
the following value in the input data: _ARGS["fields"] = "user_id,username"_. (Note that multiple column names are concated and separated by comma.)

#### 2.2.5) *Topics*

```go
    Action
    Joins       []*Join             `json:"joins,omitempty" hcl:"join,block"`
    Columns      map[string][]string `json:"columns" hcl:"columns"`
    FIELDS      string              `json:"fields,omitempty" hcl:"fields"`

    TotalForce  int    `json:"total_force,omitempty" hcl:"total_force,optional"`
    MAXPAGENO   string `json:"maxpageno,omitempty" hcl:"maxpageno,optional"`
    TOTALNO     string `json:"totalno,omitempty" hcl:"totalno,optional"`
    ROWCOUNT    string `json:"rawcount,omitempty" hcl:"rawcount,optional"`
    PAGENO      string `json:"pageno,omitempty" hcl:"pageno,optional"`
    SORTBY      string `json:"sortby,omitempty" hcl:"sortby,optional"`
    SORTREVERSE string `json:"sortreverse,omitempty" hcl:"sortreverse,optional"`
}
```

It searches many rows with pagination. The capital fields are:

Field | Default | Meaning in Input Data `ARGS`
--------- | ------- | -----------------------
_MAXPAGENO_ | "maxpageno" | how many pages in total
_TOTALNO_ | "totalno" | how many records in total
_ROWCOUNT_ | "rowcount" | how many record in each page
_PAGENO_ | "pageno" | return only data of the specific page
_SORTBY_ | "sortby" | sort the returned data by this
_SORTREVERSE_ | "sortreverse" | 1 to return the data in reverse

_TotalForce_ is defined in this way: 0 for not calculating total number of records; -1 for calculating; and 1 for optionally calculating. In the last case, if there is no input data for `ROWCOUNT` or `PAGENO`, there is no pagination information.

#### 2.2.6) *Delete*

```go
type Delete struct {
    Action
}
```

It deletes a row by the primary key.

<br />

### 2.3  *Atom*

*Atom* contains the table and actions on the table, using the *Navigate* interface:

```go
type Atom struct {
 Table
 Actions map[string]interface{} `json:"actions,omitempty" hcl:"actions,optional"`
}
```

```go
type Navigate interface {
    NonePass(action string) []string
    RunAtomContext(ctx context.Context, db *sql.DB, action string, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]map[string]interface{}, []*Nextpage, error)
}
```

where *Actions* defines an action name (as the key) and the associated *Capability* (as the value).

*NonePass* defines columns, whose values should not be passed to the next actions as constraints.

To parse _Atom_ from the json file *filename*
containing the *CRUD* capabilities and customized capabilities *custom*:

```go
func NewAtomJsonFile(filename string, custom ...map[string]Capability) (*Atom, error)
```

Note that if to write own parse function, make sure to run `Assertion` to assert right `Action` types:

```go
func (self *Atom) Assertion(custom ...map[string]Capability) error
```

<br />

### 2.4  *Nextpage*

As in MoleculeQL and gRPC, *molecule* allows action to trigger multiple actions using *Nextpage*:

```go
type Nextpage struct {
    Atom      string            `json:"atom"`                 // name of the next atom to call  
    Action     string            `json:"action"`                // action name of the next atom
    RelateItem map[string]string `json:"relateItem,omitempty"` // column name mapped to that of the next atom
    Extra      map[string]string `json:"extra,omitempty"`      // manually assign these constraints
}
```

where *Atoms* are defined as a map between names and instance of type *Atom*, and *Atom* in *Nextpage* denotes the name. *Action* is the action name, *RelateItem* the map between the current data columns to next action's columns, whose values will be used as constraints, *Extra* the manually-input constraint on the next action.

Here is a use case. There are two tables, one for family and the other for children, corresponding to atoms named `ta` and `tb` respectively.
We search the family name in `ta`, and want to show all children as well. Technically, it means we need to run a `Topics` action on *ta*. For each row returned, we run *Topics* on *tb*, constrained by the family ID in both the tables.

So *Nextpages* of *Topics* on *ta* will look like:

<details>
    <summary>Click to show the JSON string</summary>
    <p>

```json
{
    "topics" : [
        {"atom":"tb", "action":"topics", "relateItem":{"id":"id"}}
    ]
}
```

</p>
</details>

Parsing the JSON will build up a `map[string][]*Nextpage` structure.

<br />

### 2.5) Example

<details>
    <summary>Click for Example 2</summary>
    <p>



<br /><br />

## 3. `Molecule` Usage

*Molecule* describes a database

```go
type Molecule struct {
    *sql.DB
    Atoms map[string]Navigate
}
```

### 3.1 Constructor

```go
func NewMolecule(db *sql.DB, s map[string]Navigate) *Molecule
```

<br />

### 3.2 Run actions on atoms

```go
func (self *Molecule) RunContext(ctx context.Context, atom, action string, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]map[string]interface{}, error)
```

which returns the data as *[]map[string]interface{}*, and optional error.

<br />

### 3.3) Example

<details>
    <summary>Click for Example 3</summary>
    <p>

```go
package main

import (
    "log"
    "database/sql"
    "os"
    "github.com/genelet/molecule"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dbUser := os.Getenv("DBUSER")
    dbPass := os.Getenv("DBPASS")
    dbName := os.Getenv("DBNAME")
    db, err := sql.Open("mysql", dbUser + ":" + dbPass + "@/" + dbName)
    if err != nil { panic(err) }
    defer db.Close()

    db.Exec(`drop table if exists test_a`)
    db.Exec(`CREATE TABLE test_a (id int auto_increment not null primary key,
        x varchar(8), y varchar(8), z varchar(8))`)
    db.Exec(`drop table if exists test_b`)
    db.Exec(`CREATE TABLE test_b (tid int auto_increment not null primary key,
        child varchar(8), id int)`)

    ta, err := molecule.NewAtomJsonFile("test_a.json")
    if err != nil { panic(err) }
    tb, err := molecule.NewAtomJsonFile("test_b.json")
    if err != nil { panic(err) }

    molecule := &molecule.Molecule{db, map[string]molecule.Navigate{"ta":ta, "tb":tb}}

    methods := map[string]string{"LIST":"topics", "GET":"edit", "POST":"insert", "PATCH":"insupd", "PUT":"update", "DELETE":"delete"}

    var lists []map[string]interface{}
    // the 1st web requests is assumed to create id=1 to the test_a and test_b tables:
    //
    args := map[string]interface{}{"x":"a1234567","y":"b1234567","z":"temp", "child":"john"}
    if lists, err = molecule.Run("ta", methods["PATCH"], args); err != nil { panic(err) }

    // the 2nd request just updates, because [x,y] is defined to the unique in ta.
    // but create a new record to tb for id=1, since insupd triggers insert in tb
    //
    args = map[string]interface{}{"x":"a1234567","y":"b1234567","z":"zzzzz", "child":"sam"}
    if lists, err = molecule.Run("ta", methods["PATCH"], args); err != nil { panic(err) }

    // the 3rd request creates id=2
    //
    args = map[string]interface{}{"x":"c1234567","y":"d1234567","z":"e1234","child":"mary"}
    if lists, err = molecule.Run("ta", methods["POST"], args); err != nil { panic(err) }

    // the 4th request creates id=3
    //
    args = map[string]interface{}{"x":"e1234567","y":"f1234567","z":"e1234","child":"marcus"}
    if lists, err = molecule.Run("ta", methods["POST"], args); err != nil { panic(err) }

    // LIST all
    args = map[string]interface{}{}
    if lists, err = molecule.Run("ta", methods["LIST"], args); err != nil { panic(err) }
    log.Printf("LIST: %v", lists)

    // GET one
    args = map[string]interface{}{"id":1}
    if lists, err = molecule.Run("ta", methods["GET"], args); err != nil { panic(err) }
    log.Printf("GET: %v", lists)

    // DELETE
    extra := map[string]interface{}{"id":1}
    if lists, err = molecule.Run("tb", methods["DELETE"], map[string]interface{}{"tid": 1}, extra); err != nil { panic(err) }
    if lists, err = molecule.Run("tb", methods["DELETE"], map[string]interface{}{"tid": 2}, extra); err != nil { panic(err) }
    if lists, err = molecule.Run("ta", methods["DELETE"], map[string]interface{}{"id":1}); err != nil { panic(err) }

    // LIST all
    args = map[string]interface{}{}
    if lists, err = molecule.Run("ta", methods["LIST"], args); err != nil { panic(err) }
    log.Printf("LIST: %v", lists)

    os.Exit(0)
}
```

Running it will result in:

```bash
LIST: [map[id:1 ta_edit:[map[id:1 tb_topics:[map[child:john id:1 tid:1]] x:a1234567 y:b1234567 z:zzzzz]] x:a1234567 y:b1234567 z:zzzzz] map[id:2 ta_edit:[map[id:2 tb_topics:[map[child:mary id:2 tid:3]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 ta_edit:[map[id:3 tb_topics:[map[child:marcus id:3 tid:4]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]
GET: [map[id:1 tb_topics:[map[child:john id:1 tid:1]] x:a1234567 y:b1234567 z:zzzzz]]
LIST: [map[id:2 ta_edit:[map[id:2 tb_topics:[map[child:mary id:2 tid:3]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 ta_edit:[map[id:3 tb_topics:[map[child:marcus id:3 tid:4]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]
```

</p>
</details>
