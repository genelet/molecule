package godbi

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/genelet/determined/dethcl"
)

// newAtomJSONFile parse a disk file to atom
func newAtomJSONFile(fn string, custom ...Capability) (*Atom, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	atom := new(Atom)
	err = json.Unmarshal(dat, atom)
	if err != nil {
		return nil, err
	}

	atom.MergeCustomActions(custom...)
	return atom, nil
}

func TestHCLAtomMarshal(t *testing.T) {
	atom, err := newAtomJSONFile("short.json", new(SQL))
	if err != nil {
		t.Fatal(err)
	}
	bs, err := dethcl.Marshal(atom)
	if err != nil {
		t.Errorf("%s", bs)
		t.Fatal(err)
	}
}

// NewAtomHclFile parse a HCL file to atom
func newAtomHclFile(fn string, custom ...Capability) (*Atom, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	atom := new(Atom)
	err = dethcl.Unmarshal(dat, atom)
	if err != nil {
		return nil, err
	}

	atom.MergeCustomActions(custom...)
	return atom, nil
}

func TestHCLAtomUnmarshal(t *testing.T) {
	atom, err := newAtomHclFile("m_a.hcl")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := json.Marshal(atom)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != `{"tableName":"m_a","columns":[{"columnName":"x","typeName":"string","label":"x","notnull":true},{"columnName":"y","typeName":"string","label":"y","notnull":true},{"columnName":"z","typeName":"string","label":"z"},{"columnName":"id","typeName":"int","label":"id","auto":true}],"pks":["id"],"idAuto":"id","uniques":["x","y"],"actions":[{"actionName":"topics","nextpages":[{"atomName":"m_a","actionName":"edit","relateExtra":{"id":"id"}}]},{"actionName":"insert","nextpages":[{"atomName":"m_b","actionName":"insert","relateArgs":{"id":"id"}}]},{"actionName":"insupd","nextpages":[{"atomName":"m_b","actionName":"insert","relateArgs":{"id":"id"}}]},{"actionName":"edit","nextpages":[{"atomName":"m_b","actionName":"topics","relateExtra":{"id":"id"}}]},{"actionName":"update"},{"actionName":"delete"},{"actionName":"delecs"},{"actionName":"stmt","statement":""}]}` {
		t.Errorf("%s", bs)
	}
}

func TestAtomJsonParse(t *testing.T) {
	atom, err := newAtomJSONFile("m_a.json")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := json.Marshal(atom)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != `{"atomName":"m_a","tableName":"m_a","columns":[{"columnName":"x","typeName":"string","label":"x","notnull":true},{"columnName":"y","typeName":"string","label":"y","notnull":true},{"columnName":"z","typeName":"string","label":"z"},{"columnName":"id","typeName":"int","label":"id","auto":true}],"pks":["id"],"idAuto":"id","uniques":["x","y"],"actions":[{"actionName":"insert","nextpages":[{"atomName":"m_b","actionName":"insert","relateArgs":{"id":"id"}}]},{"actionName":"update"},{"actionName":"insupd","nextpages":[{"atomName":"m_b","actionName":"insert","relateArgs":{"id":"id"}}]},{"actionName":"delete"},{"actionName":"delecs"},{"actionName":"topics","nextpages":[{"atomName":"m_a","actionName":"edit","relateExtra":{"id":"id"}}]},{"actionName":"edit","nextpages":[{"atomName":"m_b","actionName":"topics","relateExtra":{"id":"id"}}]},{"actionName":"stmt","statement":""}]}` {
		t.Errorf("%s", bs)
	}
}

type SQL struct {
	Action
	Statement string `json:"statement"`
}

func (self *SQL) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]any, extra ...map[string]any) ([]any, error) {
	lists := make([]any, 0)
	dbi := &DBI{DB: db}
	var names []any
	for _, col := range t.Columns {
		names = append(names, col.ColumnName)
	}
	err := dbi.SelectSQLContext(ctx, &lists, self.Statement, names, ARGS["bravo"])
	return lists, err
}

func TestAtom(t *testing.T) {
	atom, err := newAtomJSONFile("atom.json", new(SQL))
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
					if (i == 0 && (page.AtomName != "adv_campaign")) ||
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
			if atom.AtomName != "adv_campaign" ||
				atom.TableName != "adv_campaign" ||
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
	"uniques":["x","y"]
}`
	atom := new(Atom)
	err = json.Unmarshal([]byte(str), atom)
	if err != nil {
		t.Fatal(err)
	}

	var lists []any
	// the 1st web requests is assumed to create id=1 to the m_a table
	//
	args := map[string]any{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}
	_, err = atom.RunAtom(db, "insert", args)
	if err != nil {
		t.Fatal(err)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique
	//
	args = map[string]any{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "child": "sam"}
	_, err = atom.RunAtom(db, "insupd", args)
	if err != nil {
		t.Fatal(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]any{"x": "c1234567", "y": "d1234567", "z": "e1234", "child": "mary"}
	_, err = atom.RunAtom(db, "insert", args)
	if err != nil {
		t.Fatal(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]any{"x": "e1234567", "y": "f1234567", "z": "e1234", "child": "marcus"}
	_, err = atom.RunAtom(db, "insupd", args)
	if err != nil {
		t.Fatal(err)
	}

	// GET all
	args = map[string]any{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil {
		t.Fatal(err)
	}
	// []map[string]interface {}{map[string]interface {}{"id":1, "x":"a1234567", "y":"b1234567", "z":"zzzzz"}, map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	e1 := lists[0].(map[string]any)
	e2 := lists[2].(map[string]any)
	if len(lists) != 3 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" ||
		e2["y"].(string) != "f1234567" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]any{"id": 1}
	lists, err = atom.RunAtom(db, "edit", args)
	if err != nil {
		t.Fatal(err)
	}
	e1 = lists[0].(map[string]any)
	if len(lists) != 1 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 tb_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]]

	// DELETE
	args = map[string]any{"id": 1}
	_, err = atom.RunAtom(db, "delete", args)
	if err != nil {
		t.Fatal(err)
	}

	// GET all
	args = map[string]any{}
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
	"uniques":["x","y"]
}`
	atom := new(Atom)
	err = json.Unmarshal([]byte(str), atom)
	if err != nil {
		t.Fatal(err)
	}

	var lists []any
	// the 1st web requests is assumed to create id=1 to the m_a table
	argss := []map[string]any{{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}}
	lists, err = atom.RunAtom(db, "insert", argss)
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 || lists[0].(map[string]any)["id"].(int64) != 1 {
		t.Errorf("%#v", lists)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique
	// the 3rd request creates id=2
	// the 4th request creates id=3
	argss = []map[string]any{{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "child": "sam"}, {"x": "c1234567", "y": "d1234567", "z": "e1234", "child": "mary"}, {"x": "e1234567", "y": "f1234567", "z": "e1234", "child": "marcus"}}
	lists, err = atom.RunAtom(db, "insupd", argss)
	if err != nil {
		t.Fatal(err)
	}
	// []map[string]interface {}{map[string]interface {}{"id":1, "x":"a1234567", "y":"b1234567", "z":"zzzzz"}, map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	if len(lists) != 3 || lists[0].(map[string]any)["id"].(int64) != 1 || lists[1].(map[string]any)["id"].(int64) != 2 || lists[2].(map[string]any)["id"].(int64) != 3 || lists[2].(map[string]any)["y"] != "f1234567" {
		t.Errorf("%#v", lists)
	}

	// GET all
	args := map[string]any{}
	lists, err = atom.RunAtom(db, "topics", args)
	if err != nil {
		t.Fatal(err)
	}
	// []map[string]interface {}{map[string]interface {}{"id":1, "x":"a1234567", "y":"b1234567", "z":"zzzzz"}, map[string]interface {}{"id":2, "x":"c1234567", "y":"d1234567", "z":"e1234"}, map[string]interface {}{"id":3, "x":"e1234567", "y":"f1234567", "z":"e1234"}}
	e1 := lists[0].(map[string]any)
	e2 := lists[2].(map[string]any)
	if len(lists) != 3 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" ||
		e2["y"].(string) != "f1234567" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]any{"id": 1}
	lists, err = atom.RunAtom(db, "edit", args)
	if err != nil {
		t.Fatal(err)
	}
	e1 = lists[0].(map[string]any)
	if len(lists) != 1 ||
		e1["id"].(int) != 1 ||
		e1["z"].(string) != "zzzzz" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 tb_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]]

	// DELETE
	argss = []map[string]any{{"id": 1}}
	_, err = atom.RunAtom(db, "delete", argss)
	if err != nil {
		t.Fatal(err)
	}

	// GET all
	args = map[string]any{}
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
