package godbi

import (
	"context"
	"database/sql"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Topics struct for search multiple rows by constraints
type Topics struct {
	Action
	FIELDS string `json:"fields,omitempty" hcl:"fields,optional"`

	Totalforce  int    `json:"totalforce,omitempty" hcl:"totalforce,optional"`
	MAXPAGENO   string `json:"maxpageno,omitempty" hcl:"maxpageno,optional"`
	TOTALNO     string `json:"totalno,omitempty" hcl:"totalno,optional"`
	PAGESIZE    string `json:"pagesize,omitempty" hcl:"pagesize,optional"`
	PAGENO      string `json:"pageno,omitempty" hcl:"pageno,optional"`
	SORTBY      string `json:"sortby,omitempty" hcl:"sortby,optional"`
	SORTREVERSE string `json:"sortreverse,omitempty" hcl:"sortreverse,optional"`
}

var _ Capability = (*Topics)(nil)

func (self *Topics) setDefaultElementNames() []string {
	if self.FIELDS == "" {
		self.FIELDS = "fields"
	}
	if self.SORTBY == "" {
		self.SORTBY = "sortby"
	}
	if self.SORTREVERSE == "" {
		self.SORTREVERSE = "sortreverse"
	}
	if self.PAGESIZE == "" {
		self.PAGESIZE = "pagesize"
	}
	if self.PAGENO == "" {
		self.PAGENO = "pageno"
	}
	if self.TOTALNO == "" {
		self.TOTALNO = "totalno"
	}
	if self.MAXPAGENO == "" {
		self.MAXPAGENO = "maxpageno"
	}
	return []string{self.FIELDS, self.SORTBY, self.SORTREVERSE, self.PAGESIZE, self.PAGENO, self.TOTALNO, self.MAXPAGENO}
}

// orderString outputs the ORDER BY string using information in args
func (self *Topics) orderString(t *Table, ARGS map[string]interface{}, joints ...[]*Joint) string {
	nameSortby := self.SORTBY
	nameSortreverse := self.SORTREVERSE
	namePagesize := self.PAGESIZE
	namePageno := self.PAGENO

	column := ""
	if ARGS[nameSortby] != nil {
		column = ARGS[nameSortby].(string)
	} else if joints != nil {
		table := joints[0][0]
		if table.Sortby != "" {
			column = table.Sortby
		} else {
			name := table.TableName
			if table.Alias != "" {
				name = table.Alias
			}
			name += "."
			column = name + strings.Join(t.Pks, ", "+name)
		}
	} else {
		column = strings.Join(t.Pks, ", ")
	}

	order := "ORDER BY " + column
	if _, ok := ARGS[nameSortreverse]; ok {
		order += " DESC"
	}
	if rowInterface, ok := ARGS[namePagesize]; ok {
		pagesize := 0
		switch v := rowInterface.(type) {
		case int:
			pagesize = v
		case string:
			pagesize, _ = strconv.Atoi(v)
		default:
		}
		pageno := 1
		if pnInterface, ok := ARGS[namePageno]; ok {
			switch v := pnInterface.(type) {
			case int:
				pageno = v
			case string:
				pageno, _ = strconv.Atoi(v)
			default:
			}
		} else {
			ARGS[namePageno] = 1
		}
		order += " LIMIT " + strconv.Itoa(pagesize) + " OFFSET " + strconv.Itoa((pageno-1)*pagesize)
	}

	matched, err := regexp.MatchString("[;'\"]", order)
	if err != nil || matched {
		return ""
	}
	return order
}

func (self *Topics) pagination(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) error {
	nameTotalno := self.TOTALNO
	namePagesize := self.PAGESIZE
	namePageno := self.PAGENO
	nameMaxpageno := self.MAXPAGENO

	totalforce := self.Totalforce
	// 0 means no total calculation, this is the default i.e. no report of total number of pages
	// totalforce is not allowed to pass in args for securit reason
	// -1 means total number of pages is calculated from the database
	if totalforce == 0 || ARGS[namePagesize] == nil || ARGS[namePageno] != nil {
		return nil
	}

	nt := 0
	if totalforce < -1 { // take the absolute as the total number
		nt = int(math.Abs(float64(totalforce)))
	} else if totalforce == -1 || ARGS[nameTotalno] == nil {
		if err := t.totalHashContext(ctx, db, &nt, extra...); err != nil {
			return err
		}
	} else {
		switch v := ARGS[nameTotalno].(type) {
		case int:
			nt = v
		case string:
			nt64, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return err
			}
			nt = int(nt64)
		default:
		}
	}

	ARGS[nameTotalno] = nt
	nr := 0
	switch v := ARGS[namePagesize].(type) {
	case int:
		nr = v
	case string:
		nr64, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return err
		}
		nr = int(nr64)
	default:
	}
	ARGS[nameMaxpageno] = (nt-1)/nr + 1
	return nil
}

func (self *Topics) RunAction(db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	return self.RunActionContext(context.Background(), db, t, ARGS, extra...)
}

func (self *Topics) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	self.setDefaultElementNames()
	sql, labels := t.filterPars(ARGS, self.FIELDS, self.getAllowed())
	order := self.orderString(t, ARGS)

	err := self.pagination(ctx, db, t, ARGS, extra...)
	if err != nil {
		return nil, err
	}

	newExtra := t.byConstraint(ARGS, extra...)
	if hasValue(newExtra) {
		where, values := selectCondition(newExtra, t.TableName)
		if where != "" {
			sql += "\nWHERE " + where
		}
		if order != "" {
			sql += "\n" + order
		}
		if t.dbDriver == Postgres {
			sql = questionMarkerNumber(sql)
		}

		return getSQL(ctx, db, sql, labels, values...)
	}

	if order != "" {
		sql += "\n" + order
	}
	if t.dbDriver == Postgres {
		sql = questionMarkerNumber(sql)
	}

	return getSQL(ctx, db, sql, labels)
}
