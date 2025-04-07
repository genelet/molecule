package godbi

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/genelet/determined/dethcl"
	"github.com/genelet/determined/utils"
)

// newMoleculeHclFile parse a HCL file to atom
func newMoleculeHclFile(fn string) (*Molecule, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	m := new(Molecule)
	err = dethcl.Unmarshal(dat, m)
	return m, err
}

func TestMoleculeContext(t *testing.T) {
	ta, err := newAtomJSONFile("m_a.json")
	if err != nil {
		t.Fatal(err)
	}
	tb, err := newAtomJSONFile("m_b.json")
	if err != nil {
		t.Fatal(err)
	}
	molecule := &Molecule{Atoms: []*Atom{ta, tb}}
	MoleculeGeneral(t, molecule)
}

func TestMoleculeEasy(t *testing.T) {
	dat, err := os.ReadFile("molecule.json")
	if err != nil {
		t.Fatal(err)
	}
	m := new(Molecule)
	err = json.Unmarshal(dat, m)
	if err != nil {
		t.Fatal(err)
	}
	if m.Atoms[0].Table.TableName != "m_a" ||
		m.Atoms[1].Table.TableName != "m_b" {
		t.Errorf("%#v", m.Atoms[0])
		t.Errorf("%#v", m.Atoms[1])
	}
}

func TestHCLMoleculeSimple(t *testing.T) {
	m, err := newMoleculeHclFile("molecule.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if m.Atoms[0].Table.TableName != "m_a" ||
		m.Atoms[1].Table.TableName != "m_b" {
		t.Errorf("%#v", m.Atoms[0])
		t.Errorf("%#v", m.Atoms[1])
	}
	bs, err := dethcl.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	dat, err := os.ReadFile("molecule.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != string(dat) {
		t.Errorf("m; %d=>%s", len(bs), bs)
		t.Errorf("b: %d=>%s", len(dat), dat)
	}
}

func TestHCLMoleculeOld(t *testing.T) {
	dat, err := os.ReadFile("molecule.hcl")
	if err != nil {
		t.Fatal(err)
	}
	m := new(Molecule)
	spec, err := utils.NewStruct(
		"Molecule", map[string]any{
			"Atoms": [][2]any{
				{"Atom", map[string]any{"Actions": [][2]any{
					{"insert", map[string]any{"Nextpages": []string{"Connection"}}},
					{"update", map[string]any{"Nextpages": []string{"Connection"}}},
					{"insupd", map[string]any{"Nextpages": []string{"Connection"}}},
					{"delete", map[string]any{"Nextpages": []string{"Connection"}}},
					{"delecs", map[string]any{"Nextpages": []string{"Connection"}}},
					{"topics", map[string]any{"Nextpages": []string{"Connection"}}},
					{"edit", map[string]any{"Nextpages": []string{"Connection"}}},
					{"stmt", map[string]any{"Nextpages": []string{"Connection"}}}}},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	ref := map[string]any{
		"Connection": new(Connection),
		"insupd":     new(Insupd),
		"update":     new(Update),
		"edit":       new(Edit),
		"insert":     new(Insert),
		"topics":     new(Topics),
		"delete":     new(Delete),
		"delecs":     new(Delecs),
		"stmt":       new(Stmt),
		"Atom":       new(Atom)}
	err = dethcl.UnmarshalSpec(dat, m, spec, ref)
	if err != nil {
		t.Fatal(err)
	}
	if m.Atoms[0].Table.TableName != "m_a" ||
		m.Atoms[1].Table.TableName != "m_b" {
		t.Errorf("%#v", m.Atoms[0])
		t.Errorf("%#v", m.Atoms[1])
	}
	bs, err := dethcl.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != len(dat) {
		t.Errorf("%d=>%s", len(dat), dat)
		t.Errorf("%d=>%s", len(bs), bs)
	}
}

// newMoleculeJson parses a JSON file into Molecule
func newMoleculeJSONFile(fn string) (*Molecule, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	m := new(Molecule)
	err = json.Unmarshal(dat, m)
	return m, err
}

func TestMoleculeParse(t *testing.T) {
	molecule, err := newMoleculeJSONFile("molecule.json")
	if err != nil {
		t.Fatal(err)
	}
	MoleculeGeneral(t, molecule)
}

func TestMoleculeDelecs(t *testing.T) {
	molecule, err := newMoleculeJSONFile("molecule2.json")
	if err != nil {
		t.Fatal(err)
	}
	molecule.SetLogger(DevelopLogger(os.Stdout, slog.LevelDebug))
	db, ctx, METHODS := local2Vars()
	var lists []any

	// the 1st web requests is assumed to create id=1 to the m_a and m_b tables:
	//
	args := map[string]any{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john"}
	data2 := []map[string]any{{"child": "john"}, {"child": "john2"}}
	molecule.Initialize(map[string]any{
		"m_a": map[string]any{"insupd": args},
		"m_b": map[string]any{"insert": data2},
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
		"m_b": map[string]any{"insert": data},
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

	molecule2Check(ctx, db, molecule, METHODS, t)
}

func TestMoleculeDelecs2(t *testing.T) {
	molecule, err := newMoleculeJSONFile("molecule21.json")
	if err != nil {
		t.Fatal(err)
	}
	db, ctx, METHODS := local2Vars()
	var lists []any

	// the 1st web requests is assumed to create id=1 to the m_a and m_b tables:
	//
	args := map[string]any{"x": "a1234567", "y": "b1234567", "z": "temp", "child": "john", "m_b": []map[string]any{{"child": "john"}, {"child": "john2"}}}
	//molecule.Initialize(map[string]any{
	//"m_a":map[string]any{"insupd": args},
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
	args = map[string]any{"x": "a1234567", "y": "b1234567", "z": "zzzzz", "m_b": map[string]any{"child": "sam"}}
	//raph.Initialize(map[string]any{
	//"m_a":map[string]any{"insupd": args},
	//, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["PATCH"], args); err != nil {
		panic(err)
	}

	// the 3rd request creates id=2
	//
	args = map[string]any{"x": "c1234567", "y": "d1234567", "z": "e1234", "m_b": map[string]any{"child": "mary"}}
	//raph.Initialize(map[string]any{
	//"m_a":map[string]any{"insert": args},
	//, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args); err != nil {
		panic(err)
	}

	// the 4th request creates id=3
	//
	args = map[string]any{"x": "e1234567", "y": "f1234567", "z": "e1234", "m_b": map[string]any{"child": "marcus"}}
	//raph.Initialize(map[string]any{
	//"m_a":map[string]any{"insert": args},
	//, nil)
	if _, err = molecule.RunContext(ctx, db, "m_a", METHODS["POST"], args); err != nil {
		panic(err)
	}

	molecule2Check(ctx, db, molecule, METHODS, t)
}

func TestMoleculeThreeTables(t *testing.T) {
	molecule, err := newMoleculeJSONFile("molecule3.json")
	if err != nil {
		t.Fatal(err)
	}
	MoleculeThreeGeneral(molecule, t)
}

func TestMoleculeThreeTables2(t *testing.T) {
	molecule, err := newMoleculeJSONFile("molecule31.json")
	if err != nil {
		t.Fatal(err)
	}
	MoleculeThreeGeneral(molecule, t)
}
