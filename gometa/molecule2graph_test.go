package gometa

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/genelet/molecule/godbi"
	"google.golang.org/protobuf/proto"
)

func tryit(m *godbi.Molecule, t *testing.T) {
	g := MoleculeToGraph(m, nil, "gometa", "Graph_id")
	m1, oneofs := GraphToMolecule(g)
	g1 := MoleculeToGraph(m1, oneofs, "gometa", "Graph_id")
	if !proto.Equal(g, g1) {
		if g.PackageName != g1.PackageName {
			t.Errorf("%s", g.PackageName)
			t.Errorf("%s", g1.PackageName)
		}
		for i, n := range g.Nodes {
			n1 := g1.Nodes[i]
			if !proto.Equal(n, n1) {
				t.Errorf("%s", n.String())
				t.Errorf("%s", n1.String())
			}
		}
	}
}

// newAtomJSONFile parse a disk file to atom
func newAtomJSONFile(fn string, custom ...godbi.Capability) (*godbi.Atom, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	atom := new(godbi.Atom)
	err = json.Unmarshal(dat, atom)
	if err != nil {
		return nil, err
	}

	atom.MergeCustomActions(custom...)
	return atom, nil
}

// newMoleculeJson parses a JSON file into Molecule
func newMoleculeJSONFile(fn string) (*godbi.Molecule, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	m := new(godbi.Molecule)
	err = json.Unmarshal(dat, m)
	return m, err
}

func TestMolecule2Graph(t *testing.T) {
	ta, err := newAtomJSONFile("m_a.json")
	if err != nil {
		t.Fatal(err)
	}
	tb, err := newAtomJSONFile("m_b.json")
	if err != nil {
		t.Fatal(err)
	}
	molecule := &godbi.Molecule{Atoms: []*godbi.Atom{ta, tb}}
	tryit(molecule, t)

	for _, fn := range []string{"molecule.json", "molecule2.json", "molecule21.json", "molecule3.json", "molecule31.json"} {
		molecule, err = newMoleculeJSONFile(fn)
		if err != nil {
			t.Fatal(err)
		}
		tryit(molecule, t)
	}
}
