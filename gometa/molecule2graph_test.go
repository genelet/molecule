package gometa

import (
	"testing"

	"github.com/genelet/molecule/godbi"
	"github.com/golang/protobuf/proto"
)

func tryit(m *godbi.Molecule, t *testing.T) {
	g := MoleculeToGraph(m, nil, "gometa", "Graph_id")
	m1, oneofs := GraphToMolecule(g)
	g1 := MoleculeToGraph(m1, oneofs, "gometa", "Graph_id")
	if !proto.Equal(g, g1) {
		if g.PackageName != g1.PackageName {
			t.Errorf("%s",  g.PackageName)
			t.Errorf("%s", g1.PackageName)
		}
		if g.DatabaseName != g1.DatabaseName {
			t.Errorf("%s",  g.DatabaseName)
			t.Errorf("%s", g1.DatabaseName)
		}
		for i, n := range g.Nodes {
			n1 := g1.Nodes[i]
			if !proto.Equal(n, n1) {
				t.Errorf("%s",  n.String())
				t.Errorf("%s", n1.String())
			}
		}
	}
}

func TestMolecule2Graph(t *testing.T) {
	ta, err := godbi.NewAtomJsonFile("m_a.json")
	if err != nil {
		t.Fatal(err)
	}
	tb, err := godbi.NewAtomJsonFile("m_b.json")
	if err != nil {
		t.Fatal(err)
	}
	molecule := &godbi.Molecule{Atoms:[]godbi.Navigate{ta, tb}}
	tryit(molecule, t)

	for _, fn := range []string{"molecule.json", "molecule2.json", "molecule21.json", "molecule3.json", "molecule31.json"} {
		molecule, err = godbi.NewMoleculeJsonFile(fn)
		if err != nil { t.Fatal(err) }
		tryit(molecule, t)
	}
}
