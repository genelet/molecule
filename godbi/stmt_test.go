package godbi

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestStmt(t *testing.T) {
	atom, err := newAtomJsonFile("stmt.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range atom.Actions {
		k := v.GetActionName()
		switch k {
		case "stmt":
			stmt := v.(*Stmt)
			if atom.TableName != "adv_campaign" ||
				atom.Pks[0] != "campaign_id" ||
				stmt.Labels[0].([]interface{})[0].(string) != "id" ||
				stmt.Statement != "SELECT adv.adv_id, campaign_name FROM adv_campain INNER JOIN adv USING adv_id WHERE adv_email=?" {
				t.Errorf("%#v", stmt)
			}
		default:
		}
	}
}

func TestStmtRun(t *testing.T) {
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
			"actionName": "stmt",
			"pars": ["x"],
			"labels": [["id","int"], ["name","string"]],
			"statement":"SELECT id, concat(x,y,z) FROM m_a WHERE x=?"
		}
	]
}`
	atom := new(Atom)
	err = json.Unmarshal([]byte(str), atom)
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

	// GET stmt
	args = map[string]interface{}{"x": "e1234567"}
	lists, err = atom.RunAtom(db, "stmt", args)
	if err != nil {
		t.Fatal(err)
	}
	e1 := lists[0].(map[string]interface{})
	if len(lists) != 1 ||
		e1["id"].(int) != 3 ||
		e1["name"].(string) != "e1234567f1234567e1234" {
		t.Errorf("%v", lists)
	}
	// [map[id:3 name:e1234567f1234567e1234]]

	db.Exec(`drop table if exists m_a`)
	db.Exec(`drop table if exists m_b`)
}
