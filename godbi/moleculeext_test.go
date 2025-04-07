package godbi

import (
	"context"
	"database/sql"
	"testing"
)

func local2Vars() (*sql.DB, context.Context, map[string]string) {
	db, err := getdb()
	if err != nil {
		panic(err)
	}
	db.Exec(`drop table if exists m_b`)
	db.Exec(`drop table if exists m_a`)
	db.Exec(`CREATE TABLE m_a (id int auto_increment not null primary key, x varchar(8), y varchar(8), z varchar(8))`)
	db.Exec(`CREATE TABLE m_b (tid int auto_increment not null primary key, child varchar(8), id int)`)
	return db, context.Background(), map[string]string{"LIST": "topics", "GET": "edit", "POST": "insert", "PUT": "update", "PATCH": "insupd", "DELETE": "delete"}
}

func molecule2Check(ctx context.Context, db *sql.DB, molecule *Molecule, METHODS map[string]string, t *testing.T) {
	// GET all
	args := map[string]any{}
	lists, err := molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	// [map[id:1 m_a_edit:[map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]] x:a1234567 y:b1234567 z:zzzzz] map[id:2 m_a_edit:[map[id:2 m_b_topics:[map[child:mary id:2 tid:3]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_a_edit:[map[id:3 m_b_topics:[map[child:marcus id:3 tid:4]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]
	e1 := lists[0].(map[string]any)["m_a_edit"].([]any)
	e2 := e1[0].(map[string]any)["m_b_topics"].([]any)
	if e2[0].(map[string]any)["child"].(string) != "john" || e2[1].(map[string]any)["child"].(string) != "john2" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]any{"id": 1}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["GET"], args)
	if err != nil {
		panic(err)
	}
	e2 = lists[0].(map[string]any)["m_b_topics"].([]any)
	if e2[0].(map[string]any)["child"].(string) != "john" || e2[1].(map[string]any)["child"].(string) != "john2" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:john2 id:1 tid:2] map[child:sam id:1 tid:3]] x:a1234567 y:b1234567 z:zzzzz]]

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	if len(lists) != 5 {
		t.Errorf("%v", lists)
	}

	// DELETE
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["DELETE"], map[string]any{"id": 1}); err != nil {
		panic(err)
	}

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	e1 = lists[0].(map[string]any)["m_a_edit"].([]any)
	e2 = e1[0].(map[string]any)["m_b_topics"].([]any)
	if e2[0].(map[string]any)["child"].(string) != "mary" {
		t.Errorf("%v", lists)
	}
	// [map[id:2 m_a_edit:[map[id:2 m_b_topics:[map[child:mary id:2 tid:3]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_a_edit:[map[id:3 m_b_topics:[map[child:marcus id:3 tid:4]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]

	// GET all m_b
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	if len(lists) != 2 {
		t.Errorf("%v", lists)
	}

	db.Exec(`drop table if exists m_a`)
	db.Exec(`drop table if exists m_b`)
}

func MoleculeGeneral(t *testing.T, molecule *Molecule) {
	db, ctx, METHODS := local2Vars()
	var lists []any

	// the 1st web requests is assumed to create id=1 to the m_a and m_b tables:
	//
	args := map[string]any{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}
	data := map[string]any{"child": "john"}
	orig := map[string]any{"insert": data}
	molecule.Initialize(map[string]any{"m_b": orig}, nil)
	_, err := molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args)
	if err != nil {
		panic(err)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique in ta.
	// but create a new record to tb for id=1, since insupd triggers insert in tb
	//
	args = map[string]any{"x": "a1234567", "y": "b1234567", "z": "zzzzz"}
	data = map[string]any{"child": "sam"}
	orig = map[string]any{"insert": data}
	molecule.Initialize(map[string]any{"m_b": orig}, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args); err != nil {
		panic(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]any{"x": "c1234567", "y": "d1234567", "z": "e1234"}
	data = map[string]any{"child": "mary"}
	orig = map[string]any{"insert": data}
	molecule.Initialize(map[string]any{"m_b": orig}, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args); err != nil {
		panic(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]any{"x": "e1234567", "y": "f1234567", "z": "e1234"}
	data = map[string]any{"child": "marcus"}
	orig = map[string]any{"insert": data}
	molecule.Initialize(map[string]any{"m_b": orig}, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args); err != nil {
		panic(err)
	}

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	// [map[id:1 m_a_edit:[map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]] x:a1234567 y:b1234567 z:zzzzz] map[id:2 m_a_edit:[map[id:2 m_b_topics:[map[child:mary id:2 tid:3]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_a_edit:[map[id:3 m_b_topics:[map[child:marcus id:3 tid:4]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]
	e1 := lists[0].(map[string]any)["m_a_edit"].([]any)
	e2 := e1[0].(map[string]any)["m_b_topics"].([]any)
	if e2[0].(map[string]any)["child"].(string) != "john" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]any{"id": 1}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["GET"], args)
	if err != nil {
		panic(err)
	}
	e2 = lists[0].(map[string]any)["m_b_topics"].([]any)
	if e2[0].(map[string]any)["child"].(string) != "john" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:sam id:1 tid:2]] x:a1234567 y:b1234567 z:zzzzz]]

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	if len(lists) != 4 {
		t.Errorf("%v", lists)
	}

	// DELETE
	extra := map[string]any{"id": 1}
	if _, err = molecule.RunContext(ctx, db, "m_b", METHODS["DELETE"], map[string]any{"tid": 1}, extra); err != nil {
		panic(err)
	}
	if _, err = molecule.RunContext(ctx, db, "m_b", METHODS["DELETE"], map[string]any{"tid": 2}, extra); err != nil {
		panic(err)
	}
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["DELETE"], map[string]any{"id": 1}); err != nil {
		panic(err)
	}

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	e1 = lists[0].(map[string]any)["m_a_edit"].([]any)
	e2 = e1[0].(map[string]any)["m_b_topics"].([]any)
	if e2[0].(map[string]any)["child"].(string) != "mary" {
		t.Errorf("%v", lists)
	}
	// [map[id:2 m_a_edit:[map[id:2 m_b_topics:[map[child:mary id:2 tid:3]] x:c1234567 y:d1234567 z:e1234]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_a_edit:[map[id:3 m_b_topics:[map[child:marcus id:3 tid:4]] x:e1234567 y:f1234567 z:e1234]] x:e1234567 y:f1234567 z:e1234]]

	// GET all m_b
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	if len(lists) != 2 {
		t.Errorf("%v", lists)
	}

	db.Exec(`drop table if exists m_a`)
	db.Exec(`drop table if exists m_b`)
}

func MoleculeThreeGeneral(molecule *Molecule, t *testing.T) {
	db, err := getdb()
	if err != nil {
		panic(err)
	}
	db.Exec(`drop table if exists m_b`)
	db.Exec(`drop table if exists m_ab`)
	db.Exec(`drop table if exists m_a`)
	db.Exec(`CREATE TABLE m_a (id int auto_increment not null primary key, x varchar(8), y varchar(8), z varchar(8))`)
	db.Exec(`CREATE TABLE m_ab (abid int auto_increment not null primary key, id int, tid int)`)
	db.Exec(`CREATE TABLE m_b (tid int auto_increment not null primary key, child varchar(8))`)
	ctx := context.Background()
	METHODS := map[string]string{"LIST": "topics", "GET": "edit", "POST": "insert", "PUT": "update", "PATCH": "insupd", "DELETE": "delete"}

	var lists []any
	// the 1st web requests is assumed to create id=1 to the m_a and m_b tables:
	//
	args := map[string]any{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}
	data2 := []map[string]any{{"child": "john"}, {"child": "john2"}}
	molecule.Initialize(map[string]any{
		"m_a": map[string]any{"insupd": args},
		"m_b": map[string]any{"insupd": data2},
	}, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"]); err != nil {
		panic(err)
	}
	if len(lists) != 1 {
		t.Errorf("%v", lists)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique in ta.
	// but create a new record to tb for id=1, since insupd triggers insert in tb
	//
	args = map[string]any{"x": "a1234567", "y": "b1234567", "z": "zzzzz"}
	data := map[string]any{"child": "sam"}
	molecule.Initialize(map[string]any{
		"m_a": map[string]any{"insupd": args},
		"m_b": map[string]any{"insupd": data},
	}, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"]); err != nil {
		panic(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]any{"x": "c1234567", "y": "d1234567", "z": "e1234"}
	data = map[string]any{"child": "mary"}
	molecule.Initialize(map[string]any{
		"m_a": map[string]any{"insert": args},
		"m_b": map[string]any{"insert": data},
	}, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"]); err != nil {
		panic(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]any{"x": "e1234567", "y": "f1234567", "z": "e1234"}
	data = map[string]any{"child": "marcus"}
	molecule.Initialize(map[string]any{
		"m_a": map[string]any{"insert": args},
		"m_b": map[string]any{"insert": data},
	}, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"]); err != nil {
		panic(err)
	}

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	//	t.Errorf("%v", lists)
	e1 := lists[0].(map[string]any)["m_ab_topics"].([]any)
	e21 := e1[0].(map[string]any)["m_b_topics"].([]any)
	e22 := e1[1].(map[string]any)["m_b_topics"].([]any)
	if e21[0].(map[string]any)["child"].(string) != "john" || e22[0].(map[string]any)["child"].(string) != "john2" {
		t.Errorf("%v", lists)
	}

	// GET one
	args = map[string]any{"id": 1}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["GET"], args)
	if err != nil {
		panic(err)
	}
	e1 = lists[0].(map[string]any)["m_ab_topics"].([]any)
	e21 = e1[0].(map[string]any)["m_b_topics"].([]any)
	e22 = e1[1].(map[string]any)["m_b_topics"].([]any)
	if e21[0].(map[string]any)["child"].(string) != "john" || e22[0].(map[string]any)["child"].(string) != "john2" {
		t.Errorf("%v", lists)
	}
	// [map[id:1 m_b_topics:[map[child:john id:1 tid:1] map[child:john2 id:1 tid:2] map[child:sam id:1 tid:3]] x:a1234567 y:b1234567 z:zzzzz]]

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	if len(lists) != 5 {
		t.Errorf("%v", lists)
	}

	// DELETE
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["DELETE"], map[string]any{"id": 1}); err != nil {
		panic(err)
	}

	// GET all
	args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["LIST"], args)
	if err != nil {
		panic(err)
	}
	e1 = lists[0].(map[string]any)["m_ab_topics"].([]any)
	e21 = e1[0].(map[string]any)["m_b_topics"].([]any)
	if e21[0].(map[string]any)["child"].(string) != "mary" {
		t.Errorf("%v", lists)
	}
	//[map[id:2 m_ab_topics:[map[abid:4 id:2 m_b_topics:[map[child:mary tid:4]] tid:4]] x:c1234567 y:d1234567 z:e1234] map[id:3 m_ab_topics:[map[abid:5 id:3 m_b_topics:[map[child:marcus tid:5]] tid:5]] x:e1234567 y:f1234567 z:e1234]]

	// GET all m_ab
	lists, err = molecule.RunContext(ctx, db, "m_ab", METHODS["LIST"])
	if err != nil {
		panic(err)
	}
	if len(lists) != 2 {
		t.Errorf("%v", lists)
	}

	// GET all m_b
	// args = map[string]any{}
	lists, err = molecule.RunContext(ctx, db, "m_b", METHODS["LIST"])
	if err != nil {
		panic(err)
	}
	if len(lists) != 2 {
		t.Errorf("%v", lists)
	}

	db.Exec(`drop table if exists m_b`)
	db.Exec(`drop table if exists m_ab`)
	db.Exec(`drop table if exists m_a`)
}
