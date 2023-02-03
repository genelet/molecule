package protein

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/genelet/molecule/godbi"
)

type Morph struct {
	// the field name in ARGS or extra
	Name   string   `json:"name" hcl:"name,label"`
	Fn     string   `json:"fn" hcl:"fn"`
	Pars   []string `json:"pars,omitempty" hcl:"pars,optional"`
	Fields []string `json:"fields,omitempty" hcl:"fields,optional"`
}

func (self *Morph) RunMap(db *sql.DB, args map[string]interface{}, ref map[string]func(...interface{}) (interface{}, error)) (interface{}, error) {
	fun, ok := ref[self.Fn]
	if !ok {
		return nil, nil
	}
	newArray := []interface{}{}
	if db != nil {
		newArray = []interface{}{db}
	}
	if self.Pars != nil {
		for _, par := range self.Pars {
			newArray = append(newArray, par)
		}
	}
	if self.Fields != nil {
		for _, par := range self.Fields {
			newArray = append(newArray, args[par])
		}
	}
	return fun(newArray...)
}

func RunLoop(db *sql.DB, items []*Morph, args map[string]interface{}, target map[string]interface{}) error {
	if items == nil {
		return nil
	}

	funcHash := map[string]func(...interface{}) (interface{}, error){
		"identity":   Identity,
		"identities": Identities,
		"sha1":       SHA1,
		"trimleft":   Trimleft,
		"trimright":  Trimright,
		"timestamp":  Timestamp,
	}
	funcDBHash := map[string]func(...interface{}) (interface{}, error){
		"existing": IsUnique,
		"randomid": RandomID,
	}
	ref := funcHash
	if db != nil {
		ref = funcDBHash
	}

	for _, item := range items {
		v, err := item.RunMap(db, args, ref)
		if err != nil {
			return err
		}
		if v != nil {
			target[item.Name] = v // the name is used as the field name
		}
	}
	return nil
}

func Identity(vs ...interface{}) (interface{}, error) {
	if vs == nil || len(vs) == 0 {
		return nil, fmt.Errorf("no input")
	}
	return vs[0], nil
}

func Identities(vs ...interface{}) (interface{}, error) {
	if vs == nil || len(vs) == 0 {
		return nil, fmt.Errorf("no input")
	}
	lists := make([]string, 0)
	for _, v := range vs {
		lists = append(lists, fmt.Sprintf("%v", v))
	}
	return lists, nil
}

func SHA1(vs ...interface{}) (interface{}, error) {
	strs, err := Identities(vs...)
	if err != nil {
		return nil, err
	}
	h := sha1.New()
	h.Write([]byte(strings.Join(strs.([]string), "")))
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Trimleft(vs ...interface{}) (interface{}, error) {
	if vs == nil || len(vs) <= 1 {
		return nil, fmt.Errorf("wrong input")
	}
	v0 := vs[0].(string)
	v1 := vs[1].(int)
	return v0[v1:], nil
}

func Trimright(vs ...interface{}) (interface{}, error) {
	if vs == nil || len(vs) <= 1 {
		return nil, fmt.Errorf("wrong input")
	}
	v0 := vs[0].(string)
	v1 := vs[1].(int)
	return v0[:len(v0)-v1], nil
}

func Timestamp(vs ...interface{}) (interface{}, error) {
	var t time.Time
	if vs == nil || len(vs) == 0 {
		t = time.Now()
	} else {
		i64, err := strconv.ParseInt(vs[0].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		t = time.Unix(i64, 0)
	}
	return t.Format("2006-01-02 15:04:05"), nil
}

func existing(db *sql.DB, table, field string, val interface{}) (bool, error) {
	var id int
	err := db.QueryRow("SELECT 1 FROM "+table+" WHERE "+field+"=?", val).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
	}
	return true, nil
}

func IsUnique(vs ...interface{}) (interface{}, error) {
	if vs == nil || len(vs) == 4 {
		return false, fmt.Errorf("wrong input")
	}
	db := vs[0].(*sql.DB)
	table := vs[1].(string)
	field := vs[2].(string)
	ok, err := existing(db, table, field, vs[3])
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, fmt.Errorf("%v already exists in %s", vs[3], table)
	}

	return nil, nil
}

func RandomID(vs ...interface{}) (interface{}, error) {
	n := len(vs)
	if vs == nil || (n != 3 && n != 6) {
		return nil, fmt.Errorf("wrong input")
	}
	db := vs[0].(*sql.DB)
	table := vs[1].(string)
	field := vs[2].(string)
	min := 0
	max := 2147483647
	trials := 10
	if n == 6 {
		min = vs[3].(int)
		max = vs[4].(int)
		trials = vs[5].(int)
	}

	for i := 0; i < trials; i++ {
		val := min + int(rand.Float32()*float32(max-min))
		ok, err := existing(db, table, field, val)
		if err != nil {
			return nil, err
		}
		if ok {
			continue
		}
		return val, nil
	}

	return nil, fmt.Errorf("can't get random id for %s in %s", field, table)
}

func SQLModel(vs ...interface{}) (interface{}, error) {
	if vs == nil || len(vs) < 2 {
		return false, fmt.Errorf("wrong input")
	}
	dbi := &godbi.DBI{DB: vs[0].(*sql.DB)}
	sql := vs[1].(string)
	vs = vs[2:]
	lists := make([]interface{}, 0)
	err := dbi.Select(&lists, sql, vs...)
	return lists, err
}
