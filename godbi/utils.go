package godbi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reQuestion = regexp.MustCompile(`"[^"]*"|'[^']*(?:''[^']*)*'|\?`)

func questionMarkerNumber(query string) string {
	i := 1
	repl := func(in string) string {
		if in == "?" {
			x := `$` + strconv.Itoa(i)
			i++
			return x
		}
		return in
	}
	return reQuestion.ReplaceAllStringFunc(query, repl)
}

func hasValue(extra any) bool {
	if extra == nil {
		return false
	}
	switch v := extra.(type) {
	case []string:
		if len(v) == 0 {
			return false
		}
		return hasValue(v[0])
	case []any:
		if len(v) == 0 {
			return false
		}
		return hasValue(v[0])
	case []*Joint:
		if len(v) == 0 {
			return false
		}
		return hasValue(v[0])
	case map[string]string:
		if len(v) == 0 {
			return false
		}
	case map[string]any:
		if len(v) == 0 {
			return false
		}
	case []map[string]any:
		if len(v) == 0 {
			return false
		}
		return hasValue(v[0])
	default:
	}
	return true
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if !strings.ContainsRune(chr, r) {
			return r
		}
		return -1
	}, str)
}

/*

func filtering(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func mapping(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

*/

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func grep(vs []string, t string) bool {
	return index(vs, t) >= 0
}

// return 1st bool: yes, the maps are identical
// 2nd bool: no, the maps are completely different (no common key)
func compareMap(extra, item map[string]any) (bool, bool) {
	if extra == nil || item == nil {
		return false, false
	}

	keyFound := false
	identical := true
	for k, v := range item {
		vstring := fmt.Sprintf("%v", v)
		if u, ok := extra[k]; ok {
			keyFound = true
			if fmt.Sprintf("%v", u) != vstring {
				identical = false
			}
		} else {
			identical = false
		}
	}

	return identical, keyFound
}

func grepMap(lists []map[string]any, item map[string]any) bool {
	for _, each := range lists {
		if identical, _ := compareMap(each, item); identical {
			return true
		}
	}
	return false
}

// cloneMap clones extra to a new extra
func cloneMap(hash map[string]any) map[string]any {
	if hash == nil {
		return nil
	}
	newHash := map[string]any{}
	for k, v := range hash {
		newHash[k] = v
	}
	return newHash
}

// mergeMap merges two maps
func mergeMap(extra, item map[string]any) map[string]any {
	if extra == nil {
		return item
	} else if item == nil {
		return extra
	}
	newExtra := cloneMap(extra)
	for k, v := range item {
		newExtra[k] = v
	}
	return newExtra
}

// mergeMapOr merges two maps
func mergeMapOr(extra, item map[string]any) any {
	if extra == nil {
		return item
	} else if item == nil {
		return extra
	}

	identical, keyFound := compareMap(extra, item)

	if identical {
		return cloneMap(extra)
	}

	if keyFound {
		return []map[string]any{extra, item}
	}

	// all keys are different
	return mergeMap(extra, item)
}
