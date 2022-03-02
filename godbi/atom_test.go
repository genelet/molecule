package molecule

import (
	"context"
	"database/sql"
	"testing"
)

type SQL struct {
	Action
	Statement string   `json:"statement"`
}

func (self *SQL) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	lists := make([]interface{}, 0)
	dbi := &DBI{DB: db}
	var names []interface{}
	for _, col := range t.Columns {
		names = append(names, col.ColumnName)
	}
	err := dbi.SelectSQLContext(ctx, &lists, self.Statement, names, ARGS["bravo"])
	return lists, err
}

func TestAtom(t *testing.T) {
	custom := new(SQL)
	custom.ActionName = "sql"
	atom, err := NewAtomJsonFile("atom.json", custom)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range atom.Actions {
		k := v.GetActionName()
		switch k {
		case "topics":
			topics := v.(*Topics)
			if topics.Nextpages != nil {
				for i, page := range topics.Nextpages {
					if (i == 0 && (page.TableName != "adv_campaign")) ||
						(i == 1 && (page.RelateExtra["campaign_id"] != "campaign_id")) {
						t.Errorf("%#v", page)
					}
				}
			}
		case "update":
			update := v.(*Update)
			if update.Empties[0] != "created" {
				t.Errorf("%#v", update)
			}
		case "sql":
			sql := v.(*SQL)
			if atom.TableName != "adv_campaign" ||
				atom.Pks[0] != "campaign_id" ||
				sql.Nextpages[0].ActionName != "topics" ||
				sql.Statement != "SELECT x, y, z FROM a WHERE b=?" {
				t.Errorf("%#v", sql)
			}
		default:
		}
	}
}

func TestAtomRun(t *testing.T) {
	db, err := getdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`drop table if exists m_a`)
	db.Exec(`CREATE TABLE m_a (id int auto_increment not null primary key,
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
	{
		"isDo":true,
		"actionName": "insert"
	},
	{
		"isDo":true,
		"actionName": "insupd"
	},
	{
		"actionName": "delete"
	},
	{
		"actionName": "topics"
	},
	{
		"actionName": "edit"
	}
]}`
	atom, err := NewAtomJson([]byte(str))
	if err != nil {
		t.Fatal(err)
	}

	var lists []interface{}
	// the 1st web requests is assumed to create id=1 to the m_a table
	//
	args := map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}
	lists, err = atom.RunAtom(db, "insert", args)
	if err != nil {
		t.Fatal(err)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique
	//
	args = map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "child": "sam"}
	lists, err = atom.RunAtom(db, "insupd", args)
	if err != nil {
		t.Fatal(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]interface{}{"x": "c1234567", "y": "d1234567", "z": "e1234", "child": "mary"}
	lists, err = atom.RunAtom(db, "insert", args)
	if err != nil {
		t.Fatal(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]interface{}{"x": "e1234567", "y": "f1234567", "z": "e1234", "child": "marcus"}
	lists, err = atom.RunAtom(db, "insupd", args)
	if err != nil {
		t.Fatal(err)
	}

	// GET all
	args = map[string]interface{}{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil {
		t.Fatal(err)
	}
	// []map[string]interface {}{map[string]interface {}{"id":1, "x":"a1234567", "y":"b1234567", "z":"zzzzz"}, map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	e1 := lists[0].(map[string]interface{})
	e2 := lists[2].(map[string]interface{})
	if len(lists) != 3 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" ||
		e2["y"].(string) != "f1234567" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]interface{}{"id": 1}
	lists, err = atom.RunAtom(db, "edit", args)
	if err != nil {
		t.Fatal(err)
	}
	e1 = lists[0].(map[string]interface{})
	if len(lists) != 1 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 tb_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]]

	// DELETE
	args = map[string]interface{}{"id": 1}
	lists, err = atom.RunAtom(db, "delete", args)
	if err != nil {
		t.Fatal(err)
	}

	// GET all
	args = map[string]interface{}{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil {
		t.Fatal(err)
	}
	//[]map[string]interface {}{map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	if len(lists) != 2 {
		t.Errorf("%v", lists)
	}

	db.Exec(`drop table if exists m_a`)
	db.Exec(`drop table if exists m_b`)
}

func TestAtomRunMultiple(t *testing.T) {
	db, err := getdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`drop table if exists m_a`)
	db.Exec(`CREATE TABLE m_a (id int auto_increment not null primary key,
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
	{
		"actionName": "insert"
	},
	{
		"actionName": "insupd"
	},
	{
		"actionName": "delete"
	},
	{
		"actionName": "topics"
	},
	{
		"actionName": "edit"
	}
]}`
	atom, err := NewAtomJson([]byte(str))
	if err != nil {
		t.Fatal(err)
	}

	var lists []interface{}
	// the 1st web requests is assumed to create id=1 to the m_a table
	argss := []map[string]interface{}{{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}}
	lists, err = atom.RunAtom(db, "insert", argss)
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 || lists[0].(map[string]interface{})["id"].(int64) != 1 {
		t.Errorf("%#v", lists)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique
	// the 3rd request creates id=2
	// the 4th request creates id=3
	argss = []map[string]interface{}{{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "child": "sam"}, {"x": "c1234567", "y": "d1234567", "z": "e1234", "child": "mary"}, {"x": "e1234567", "y": "f1234567", "z": "e1234", "child": "marcus"}}
	lists, err = atom.RunAtom(db, "insupd", argss)
	if err != nil {
		t.Fatal(err)
	}
	// []map[string]interface {}{map[string]interface {}{"id":1, "x":"a1234567", "y":"b1234567", "z":"zzzzz"}, map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	if len(lists) != 3 || lists[0].(map[string]interface{})["id"].(int64) != 1 || lists[1].(map[string]interface{})["id"].(int64) != 2 || lists[2].(map[string]interface{})["id"].(int64) != 3 || lists[2].(map[string]interface{})["y"] != "f1234567" {
		t.Errorf("%#v", lists)
	}

	// GET all
	args := map[string]interface{}{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil {
		t.Fatal(err)
	}
	// []map[string]interface {}{map[string]interface {}{"id":1, "x":"a1234567", "y":"b1234567", "z":"zzzzz"}, map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	e1 := lists[0].(map[string]interface{})
	e2 := lists[2].(map[string]interface{})
	if len(lists) != 3 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" ||
		e2["y"].(string) != "f1234567" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]interface{}{"id": 1}
	lists, err = atom.RunAtom(db, "edit", args)
	if err != nil {
		t.Fatal(err)
	}
	e1 = lists[0].(map[string]interface{})
	if len(lists) != 1 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 tb_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]]

	// DELETE
	argss= []map[string]interface{}{{"id": 1}}
	lists, err = atom.RunAtom(db, "delete", args)
	if err != nil {
		t.Fatal(err)
	}

	// GET all
	args = map[string]interface{}{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil {
		t.Fatal(err)
	}
	//[]map[string]interface {}{map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	if len(lists) != 2 {
		t.Errorf("%v", lists)
	}

	db.Exec(`drop table if exists m_a`)
	db.Exec(`drop table if exists m_b`)
}
