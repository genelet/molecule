package godbi

import (
	"encoding/json"
	"testing"

	"github.com/genelet/determined/dethcl"
)

func TestHCLConnection(t *testing.T) {
	str := `{"atomName":"adv_campaign", "actionName":"edit", "relateExtra":{"campaign_id":"c_id"}, "relateArgs":{"x":"firstname"}}`
	page := new(Connection)
	err := json.Unmarshal([]byte(str), page)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := dethcl.Marshal(page)
	if err != nil {
		t.Fatal(err)
	}
	c := new(Connection)
	err = dethcl.Unmarshal(bs, c)
	if err != nil {
		t.Fatal(err)
	}
	if c.RelateArgs["x"] != "firstname" || c.RelateExtra["campaign_id"] != "c_id" {
		t.Errorf("%#v", c)
	}
}

func TestConnection(t *testing.T) {
	str := `{"atomName":"adv_campaign", "actionName":"edit", "relateExtra":{"campaign_id":"c_id"}, "relateArgs":{"x":"firstname"}}`
	page := new(Connection)
	err := json.Unmarshal([]byte(str), page)
	if err != nil {
		t.Fatal(err)
	}

	item := map[string]interface{}{"x": "a", "campaign_id": 123}
	nextArgs := page.nextArgs(item).(map[string]interface{})
	if nextArgs["firstname"] != "a" {
		t.Errorf("%#v", nextArgs)
	}
	nextExtra := page.nextExtra(item)
	if nextExtra["c_id"] != 123 {
		t.Errorf("%#v", nextExtra)
	}

	extra := map[string]interface{}{"y": "b", "asset": "what"}
	hash := mergeMap(nextExtra, extra)
	if hash["y"].(string) != "b" ||
		hash["asset"].(string) != "what" ||
		hash["c_id"].(int) != 123 {
		t.Errorf("%#v", hash)
	}

	item = map[string]interface{}{"x": "a", "item_id": 123}
	arg := map[string]interface{}{"y": "b", "asset": "what"}
	cArg := cloneArgs(arg).(map[string]interface{})
	aArg := mergeArgs(arg, item).(map[string]interface{})
	if len(cArg) != 2 || cArg["y"] != "b" || cArg["asset"] != "what" {
		t.Errorf("%#v", cArg)
	}
	if len(aArg) != 4 || aArg["item_id"] != 123 {
		t.Errorf("%#v", aArg)
	}

	args := []map[string]interface{}{{"y": "b", "asset": "what"},
		{"y": "bb", "asset": "whatwhat", "size_id": 777}}
	cArgs := cloneArgs(args).([]interface{})
	aArgs := mergeArgs(args, item).([]interface{})
	//[]map[string]interface {}{map[string]interface {}{"asset":"what", "y":"b"}, map[string]interface {}{"asset":"whatwhat", "size_id":777, "y":"bb"}}
	if len(cArgs) != 2 || cArgs[0].(map[string]interface{})["y"] != "b" || cArgs[1].(map[string]interface{})["asset"] != "whatwhat" {
		t.Errorf("%#v", cArgs)
	}
	//[]map[string]interface {}{map[string]interface {}{"asset":"what", "item_id":123, "x":"a", "y":"b"}, map[string]interface {}{"asset":"whatwhat", "item_id":123, "size_id":777, "x":"a", "y":"bb"}}
	if len(aArgs) != 2 || len(aArgs[0].(map[string]interface{})) != 4 || aArgs[0].(map[string]interface{})["item_id"].(int) != 123 || len(aArgs[1].(map[string]interface{})) != 5 || aArgs[1].(map[string]interface{})["item_id"].(int) != 123 || aArgs[1].(map[string]interface{})["size_id"].(int) != 777 {
		t.Errorf("%#v", aArgs)
	}
}

func TestConnectionNextArgs(t *testing.T) {
	p := &Connection{AtomName: "PersonTeacher", ActionName: "insert", RelateArgs: map[string]string{"ALL": "ALL", "PersonTeacher_id": "advisors"}}
	args := map[string]interface{}{"PersonTeacher_id": 4, "advisors": []interface{}{map[string]interface{}{"endYear": 1956, "fullname": "name 31", "school": "school31", "startYear": 1951}}, "endYear": 1960, "fullname": "name 22", "school": "school22", "startYear": 1957}
	na := p.nextArgs(args).(map[string]interface{})
	// map[string]interface {}{"PersonTeacher_id":4, "advisors":[]interface {}{map[string]interface {}{"endYear":1956, "fullname":"name 31", "school":"school31", "startYear":1951}}, "endYear":1960, "fullname":"name 22", "school":"school22", "startYear":1957}
	if na["PersonTeacher_id"].(int) != 4 ||
		na["fullname"].(string) != "name 22" ||
		len(na["advisors"].([]interface{})) != 1 {
		t.Errorf("%#v", na)
	}
}

func TestConnectionAppend(t *testing.T) {
	args := map[string]interface{}{"id": 1}
	lists := []map[string]interface{}{{"child": "john", "tid": 1}, {"child": "john2", "tid": 2}}
	var newArgs []interface{}
	var newExtra []map[string]interface{}
	p := &Connection{AtomName: "m_b", ActionName: "insert", RelateArgs: map[string]string{"tid": "tid"}}
	for _, item := range lists {
		newArgs = append(newArgs, mergeArgs(args, p.nextArgs(item).(map[string]interface{})))
		if p.nextExtra(item) != nil {
			newExtra = append(newExtra, p.nextExtra(item))
		}
	}
	if len(newArgs) != 2 ||
		newArgs[0].(map[string]interface{})["id"] != 1 ||
		newArgs[0].(map[string]interface{})["tid"] != 1 ||
		newArgs[1].(map[string]interface{})["id"] != 1 ||
		newArgs[1].(map[string]interface{})["tid"] != 2 {
		t.Errorf("%#v", newArgs)
	}
	if newExtra != nil {
		t.Errorf("%#v", newExtra)
	}
}
