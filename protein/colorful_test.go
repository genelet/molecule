package protein

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestColorful(t *testing.T) {
	bs, err := ioutil.ReadFile("tmp/adv/adv-adv.json")
	if err != nil {
		t.Fatal(err)
	}
	colorful := Colorful{}
	err = json.Unmarshal(bs, &colorful)
	if err != nil {
		t.Fatal(err)
	}

	if !colorful.Color.IsUser || colorful.Atom.Table.Pks[0] != "adv_id" {
		t.Errorf("%#v", colorful)
		t.Errorf("%#v", colorful.Atom)
		t.Errorf("%#v", colorful.Atom.Table.Pks[0])
	}
}
