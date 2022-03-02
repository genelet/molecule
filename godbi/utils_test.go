package godbi

import (
	"testing"
)

func TestUtils(t *testing.T) {
	str := "abcdefg_+hijk=="
	newStr := stripchars(str, "df+=")
	if "abceg_hijk" != newStr {
		t.Errorf("%s %s wanted", str, newStr)
	}
	x := []string{str, newStr, "abc"}
	if grep(x, "abcZ") {
		t.Errorf("%s wrong matched", "abcZ")
	}
	if grep(x, "abc") == false {
		t.Errorf("%s matched", "abc")
	}
	if grep([]string{"child", "tid"}, "tid") == false {
		t.Errorf("%#v does not match %s", []string{"child", "tid"}, "tid")
	}
	x1 := []interface{}{"a", "b"}
	x2 := map[string]interface{}{"a": "b"}
	x3 := make([]interface{}, 0)
	x4 := make(map[string]interface{})
	if !hasValue(x1) {
		t.Errorf("%v", x1)
	}
	if !hasValue(x2) {
		t.Errorf("%v", x2)
	}
	if hasValue(x3) {
		t.Errorf("%v", x3)
	}
	if hasValue(x4) {
		t.Errorf("%v", x4)
	}

	query := "INSERT INTO x (col1, col2) VALUE (?, ?), (?, ?)"
	marked := questionMarkerNumber(query)
	if marked != "INSERT INTO x (col1, col2) VALUE ($1, $2), ($3, $4)" {
		t.Errorf("%s=>%s", query, marked)
	}
}

func TestMapUtils(t *testing.T) {
	map1 := map[string]interface{}{"a":"x", "b":"y", "c":"z"}
	map2 := map[string]interface{}{"a":"x", "b":"y", "c":"z"}
	identical, keyFound := compareMap(map1, map2)
	if !identical {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	map2 = map[string]interface{}{"a":"x", "b":"y", "c":"zz"}
	identical, keyFound = compareMap(map1, map2)
	if identical || !keyFound {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	map2 = map[string]interface{}{"a1":"x", "b1":"y", "c1":"z"}
	identical, keyFound = compareMap(map1, map2)
	if identical || keyFound {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	lists := []map[string]interface{}{{"a":"x1", "b":"y", "c":"z"}, {"a":"x", "b":"y"}, {"a":"x3", "b":"y3", "c":"z3"}, {"a1":"x", "b1":"y", "c1":"z"}}
	got := grepMap(lists, map1)
	if got {
		t.Errorf("got: %t, lists: %v, item: %v", got, lists, map1)
	}
	got = grepMap(lists, map2)
	if !got {
		t.Errorf("got: %t, lists: %v, item: %v", got, lists, map2)
	}

	map1 = map[string]interface{}{"a":"x", "b":"y", "c":"z"}
	cmap1 := CloneMap(map1)
	identical, keyFound = compareMap(map1, cmap1)
	if !identical {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	map2 = map[string]interface{}{"a":"x1", "b":"y1", "c":"z1"}
	cmap12 := MergeMap(map1, map2) // map2 will replace all values in map1
	identical, keyFound = compareMap(cmap12, map2)
	if !identical {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	map2 = map[string]interface{}{"a1":"x1", "b1":"y1", "c1":"z1"}
	cmap12 = MergeMap(map1, map2)
	identical, keyFound = compareMap(cmap12, map[string]interface{}{"a1":"x1", "b1":"y1", "c1":"z1","a":"x", "b":"y", "c":"z"})
	if !identical {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	map2 = map[string]interface{}{"a":"x1", "b1":"y1", "c1":"z1"}
	cmap12 = MergeMap(map1, map2)
	identical, keyFound = compareMap(cmap12, map[string]interface{}{"a":"x1", "b1":"y1", "c1":"z1", "b":"y", "c":"z"})
	if !identical {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}

	map1 = map[string]interface{}{"a":"x", "b":"y", "c":"z"}
	map2 = map[string]interface{}{"a":"x", "b":"y", "c":"z1"}
	map12 := mergeMapOr(map1, map2)
	if len(map12.([]map[string]interface{})) != 2 {
		t.Errorf("%v", map12)
	}
	map2 = map[string]interface{}{"a1":"x", "b":"y", "c":"z"}
	map12 = mergeMapOr(map1, map2)
	if len(map12.([]map[string]interface{})) != 2 {
		t.Errorf("%v", map12)
	}

	map2 = map[string]interface{}{"a":"x", "b":"y", "c":"z"}
	map12 = mergeMapOr(map1, map2) // identical
	if len(map12.(map[string]interface{})) != 3 {
		t.Errorf("%v", map12)
	}
	map2 = map[string]interface{}{"a1":"x1", "b1":"y1", "c1":"z1"}
	map12 = MergeMap(map1, map2) // no common key
	identical, keyFound = compareMap(map12.(map[string]interface{}), map[string]interface{}{"a1":"x1", "b1":"y1", "c1":"z1", "a":"x", "b":"y", "c":"z"})
	if !identical {
		t.Errorf("identical: %t, keyFound: %t", identical, keyFound)
	}
}
