package godbi

import (
"encoding/json"
"io/ioutil"
	"testing"
)

func TestMoleculeContext(t *testing.T) {
	ta, err := NewAtomJsonFile("m_a.json")
	if err != nil {
		t.Fatal(err)
	}
	tb, err := NewAtomJsonFile("m_b.json")
	if err != nil {
		t.Fatal(err)
	}
	molecule := &Molecule{Atoms:[]Navigate{ta, tb}}
	MoleculeGeneral(t, molecule)
}

func TestMoleculeEasy(t *testing.T) {
	dat, err := ioutil.ReadFile("molecule.json")
	if err != nil { t.Fatal(err) }
	m := new(Molecule)
	err = json.Unmarshal(dat, m)
	if err != nil { t.Fatal(err) }
}

func TestMoleculeParse(t *testing.T) {
	molecule, err := NewMoleculeJsonFile("molecule.json")
	if err != nil {
		t.Fatal(err)
	}
	MoleculeGeneral(t, molecule)
}

func TestMoleculeDelecs(t *testing.T) {
	molecule, err := NewMoleculeJsonFile("molecule2.json")
	if err != nil {
		t.Fatal(err)
	}
	db, ctx, METHODS := local2Vars()
	var lists []interface{}

	// the 1st web requests is assumed to create id=1 to the m_a and m_b tables:
	//
	args := map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}
    data2 := []map[string]interface{}{{"child": "john"}, {"child": "john2"}}
	molecule.Initialize(map[string]interface{}{
		"m_a":map[string]interface{}{"insupd": args},
		"m_b":map[string]interface{}{"insert": data2},
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
	args = map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "zzzzz"}
    data:= map[string]interface{}{"child": "sam"}
	molecule.Initialize(map[string]interface{}{
		"m_a":map[string]interface{}{"insupd": args},
		"m_b":map[string]interface{}{"insert": data},
	}, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"]); err != nil {
		panic(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]interface{}{"x": "c1234567", "y": "d1234567", "z": "e1234"}
    data = map[string]interface{}{"child": "mary"}
	molecule.Initialize(map[string]interface{}{
		"m_a":map[string]interface{}{"insert": args},
		"m_b":map[string]interface{}{"insert": data},
	}, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"]); err != nil {
		panic(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]interface{}{"x": "e1234567", "y": "f1234567", "z": "e1234"}
    data = map[string]interface{}{"child": "marcus"}
	molecule.Initialize(map[string]interface{}{
		"m_a":map[string]interface{}{"insert": args},
		"m_b":map[string]interface{}{"insert": data},
	}, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"]); err != nil {
		panic(err)
	}

	molecule2Check(ctx, db, molecule, METHODS, t)
}

func TestMoleculeDelecs2(t *testing.T) {
	molecule, err := NewMoleculeJsonFile("molecule21.json")
	if err != nil {
		t.Fatal(err)
	}
	db, ctx, METHODS := local2Vars()
	var lists []interface{}

	// the 1st web requests is assumed to create id=1 to the m_a and m_b tables:
	//
	args := map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john", "m_b": []map[string]interface{}{{"child": "john"}, {"child": "john2"}}}
//molecule.Initialize(map[string]interface{}{
//"m_a":map[string]interface{}{"insupd": args},
//, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args); err != nil {
		panic(err)
	}
	if len(lists) != 1 {
		t.Errorf("%v", lists)
	}

	// the 2nd request just updates, becaues [x,y] is defined to the unique in ta.
	// but create a new record to tb for id=1, since insupd triggers insert in tb
	//
	args = map[string]interface{}{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "m_b": map[string]interface{}{"child": "sam"}}
//raph.Initialize(map[string]interface{}{
//"m_a":map[string]interface{}{"insupd": args},
//, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args); err != nil {
		panic(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]interface{}{"x": "c1234567", "y": "d1234567", "z": "e1234", "m_b": map[string]interface{}{"child": "mary"}}
//raph.Initialize(map[string]interface{}{
//"m_a":map[string]interface{}{"insert": args},
//, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args); err != nil {
		panic(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]interface{}{"x": "e1234567", "y": "f1234567", "z": "e1234", "m_b": map[string]interface{}{"child": "marcus"}}
//raph.Initialize(map[string]interface{}{
//"m_a":map[string]interface{}{"insert": args},
//, nil)
	if lists, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args); err != nil {
		panic(err)
	}

	molecule2Check(ctx, db, molecule, METHODS, t)
}

func TestMoleculeThreeTables(t *testing.T) {
	molecule, err := NewMoleculeJsonFile("molecule3.json")
	if err != nil {
		t.Fatal(err)
	}
	MoleculeThreeGeneral(molecule, t)
}

func TestMoleculeThreeTables2(t *testing.T) {
	molecule, err := NewMoleculeJsonFile("molecule31.json")
	if err != nil {
		t.Fatal(err)
	}
	MoleculeThreeGeneral(molecule, t)
}
