package protein

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestProteinJSON(t *testing.T) {
	protein, _, err := getProtein()
	if err != nil { t.Fatal(err) }

	name := "tmp/oneInOneProtein.json"

	f, err := os.Create(name)
	if err != nil { t.Fatal(err) }
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(protein)	
	if err != nil { t.Fatal(err) }
	
	g, err := os.Open(name)
	if err != nil { t.Fatal(err) }
	defer g.Close()
	decoder := json.NewDecoder(g)
	protein1 := &Protein{}
	err = decoder.Decode(protein1)	
	if err != nil { t.Fatal(err) }

	if protein.DBDriver != protein1.DBDriver ||
	protein.Teams["admin"].String()  != protein1.Teams["admin"].String() ||
	protein.Teams["adv"].String()    != protein1.Teams["adv"].String() ||
	protein.Teams["pub"].String()    != protein1.Teams["pub"].String() ||
	protein.Teams["public"].String() != protein1.Teams["public"].String() {
		t.Errorf("%#v", protein.Teams["adv"])
		t.Errorf("%#v", protein1.Teams["adv"])
	}
}

func TestProteinRest(t *testing.T) {
	protein, db, err := getProtein()
	if err != nil { t.Fatal(err) }
	defer db.Close()

	ctx := context.Background()
	token := "0"
	team := "adv"
	args := map[string]interface{}{"adv_id":1}
	extra := make(map[string]interface{})
	lists, err := protein.RunContext(ctx, db, token, team, "adv_campaign", "topics", args, extra)
	if err != nil { t.Fatal(err) }
	item := lists[0].(map[string]interface{})
	if item["campaign_id"].(int) != 1 || item["campaign_id_sign"].(string) != "3007526f8abed158c8c7002a854dec29caf801f2" || item["campaign_name"].(string) != "camp 001" {
		t.Errorf("%#v", item)
	}

	args = map[string]interface{}{"campaign_id":1, "campaign_id_sign":"3007526f8abed158c8c7002a854dec29caf801f2"}
	extra = make(map[string]interface{})
	lists, err = protein.RunContext(ctx, db, token, team, "adv_item", "topics", args, extra)
    if err != nil { t.Fatal(err) }
    item = lists[0].(map[string]interface{})
	if item["item_id"].(int) != 1 || item["item_id_sign"].(string) != "5bc219fff551de4d0f014be7b1c93f5224507be7" || item["item_name"].(string) != "item 001" {
		t.Errorf("%#v", item)
	}
}
