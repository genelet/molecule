package engine

/*
import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/genelet/molecule/godbi"
	_ "github.com/go-sql-driver/mysql"
)

func TestRunContext(t *testing.T) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := "summer"
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	admin := new(Team)
	err = NewTeamJsonFile(admin, "tmp", "admin.json")
	if err != nil { t.Fatal(err) }
	adv := new(Team)
	err = NewTeamJsonFile(adv, "tmp", "adv.json")
	if err != nil {
		t.Fatal(err)
	}

	args := map[string]interface{}{}
	extra := map[string]interface{}{"adv_id": 1}
	lists, err := admin.RunContext(context.Background(), db, "summer", godbi.DBTypeByName("mysql"), "", "adv_campaign", "topics", args, extra)
	if err != nil {
		t.Fatal(err)
	}
	list5 := lists[4].(map[string]interface{})
	if len(lists) != 5 ||
		list5["access_order"].(string) != "Inherit" ||
		list5["active"].(string) != "Yes" ||
		list5["adv_id"].(int) != 1 ||
		list5["campaign_id"].(int) != 5 ||
		list5["created"].(string) != "2020-07-17 00:01:03" {
		t.Errorf("%#v", list5)
	}

	args = map[string]interface{}{"adv_id": 1}
	extra = map[string]interface{}{}
	lists, err = adv.RunContext(context.Background(), db, "summer", godbi.DBTypeByName("mysql"), "token", "adv_campaign", "topics", args, extra)
	if err != nil {
		t.Fatal(err)
	}
	list5 = lists[4].(map[string]interface{})
	if len(lists) != 5 ||
		list5["access_order"].(string) != "Inherit" ||
		list5["active"].(string) != "Yes" ||
		list5["adv_id"].(int) != 1 ||
		list5["campaign_id"].(int) != 5 ||
		list5["created"].(string) != "2020-07-17 00:01:03" {
		t.Errorf("%#v", list5)
	}

	list1 := lists[0].(map[string]interface{})
	args = map[string]interface{}{
		"campaign_id":      list1["campaign_id"],
		"campaign_id_sign": list1["campaign_id_sign"]}
	extra = map[string]interface{}{}
	lists, err = adv.RunContext(context.Background(), db, "summer", godbi.DBTypeByName("mysql"), "token", "adv_item", "topics", args, extra)
	if err != nil {
		t.Fatal(err)
	}
	list5 = lists[4].(map[string]interface{})
	if len(lists) != 5 ||
		list5["item_id"].(int) != 5 ||
		list5["active"].(string) != "Yes" ||
		list5["cost"].(float64) != 5 ||
		list5["campaign_id"].(int) != 1 ||
		list5["startx"].(string) != "2020-07-17 00:01:03" {
		t.Errorf("%#v", lists)
	}
}
*/
