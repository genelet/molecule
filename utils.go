package molecule

import (
	"fmt"
	"regexp"
	"strings"
	"strconv"
)

func questionMarkerNumber(query string) string {
	re := regexp.MustCompile(`\?`)
	i := 1
	repl := func(in string) string {
		x := `$` + strconv.Itoa(i)
		i++
		return x
	}
	return re.ReplaceAllStringFunc(query, repl)
}

func hasValue(extra interface{}) bool {
	if extra == nil {
		return false
	}
	switch v := extra.(type) {
	case []string:
		if len(v) == 0 {
			return false
		}
		return hasValue(v[0])
	case []interface{}:
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
	case map[string]interface{}:
		if len(v) == 0 {
			return false
		}
	case []map[string]interface{}:
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
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

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
//
func compareMap(extra, item map[string]interface{}) (bool, bool) {
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

func grepMap(lists []map[string]interface{}, item map[string]interface{}) bool {
	for _, each := range lists {
		if identical, _ := compareMap(each, item); identical {
			return true
		}
	}
	return false
}

// CloneMap clones extra to a new extra
//
func CloneMap(hash map[string]interface{}) map[string]interface{} {
	if hash == nil {
		return nil
	}
	newHash := map[string]interface{}{}
	for k, v := range hash {
		newHash[k] = v
	}
	return newHash
}

// MergeMap merges two maps
//
func MergeMap(extra, item map[string]interface{}) map[string]interface{} {
	if extra == nil {
		return item
	} else if item == nil {
		return extra
	}
	newExtra := CloneMap(extra)
	for k, v := range item {
		newExtra[k] = v
	}
	return newExtra
}

// mergeMapOr merges two maps
//
func mergeMapOr(extra, item map[string]interface{}) interface{} {
	if extra == nil {
		return item
	} else if item == nil {
		return extra
	}

	identical, keyFound := compareMap(extra, item)

	if identical {
		return CloneMap(extra)
	}

	if keyFound {
		return []map[string]interface{}{extra, item}
	}

	// all keys are different
	return MergeMap(extra, item)
}
