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

func (t *Topics) setDefaultElementNames() []string {
	if t.FIELDS == "" {
		t.FIELDS = "fields"
	}
	if t.SORTBY == "" {
		t.SORTBY = "sortby"
	}
	if t.SORTREVERSE == "" {
		t.SORTREVERSE = "sortreverse"
	}
	if t.PAGESIZE == "" {
		t.PAGESIZE = "pagesize"
	}
	if t.PAGENO == "" {
		t.PAGENO = "pageno"
	}
	if t.TOTALNO == "" {
		t.TOTALNO = "totalno"
	}
	if t.MAXPAGENO == "" {
		t.MAXPAGENO = "maxpageno"
	}
	return []string{t.FIELDS, t.SORTBY, t.SORTREVERSE, t.PAGESIZE, t.PAGENO, t.TOTALNO, t.MAXPAGENO}
}

// orderString outputs the ORDER BY string using information in args
func (t *Topics) orderString(table *Table, args map[string]any, joints ...[]*Joint) string {
	nameSortby := t.SORTBY
	nameSortreverse := t.SORTREVERSE
	namePagesize := t.PAGESIZE
	namePageno := t.PAGENO

	column := ""
	if args[nameSortby] != nil {
		column = args[nameSortby].(string)
	} else if joints != nil {
		j := joints[0][0]
		if j.Sortby != "" {
			column = j.Sortby
		} else {
			name := j.TableName
			if j.Alias != "" {
				name = j.Alias
			}
			name += "."
			column = name + strings.Join(table.Pks, ", "+name)
		}
	} else {
		column = strings.Join(table.Pks, ", ")
	}

	order := "ORDER BY " + column
	if _, ok := args[nameSortreverse]; ok {
		order += " DESC"
	}
	if rowInterface, ok := args[namePagesize]; ok {
		pagesize := 0
		switch v := rowInterface.(type) {
		case int:
			pagesize = v
		case string:
			pagesize, _ = strconv.Atoi(v)
		default:
		}
		pageno := 1
		if pnInterface, ok := args[namePageno]; ok {
			switch v := pnInterface.(type) {
			case int:
				pageno = v
			case string:
				pageno, _ = strconv.Atoi(v)
			default:
			}
		} else {
			args[namePageno] = 1
		}
		if pagesize > 0 {
			if pageno < 1 {
				pageno = 1
			}
			order += " LIMIT " + strconv.Itoa(pagesize) + " OFFSET " + strconv.Itoa(pagesize*(pageno-1))
		}
	}

	matched, err := regexp.MatchString("[;'\"]", order)
	if err != nil || matched {
		return ""
	}
	return order
}

func (t *Topics) pagination(ctx context.Context, db *sql.DB, table *Table, args map[string]any, extra ...map[string]any) error {
	nameTotalno := t.TOTALNO
	namePagesize := t.PAGESIZE
	namePageno := t.PAGENO
	nameMaxpageno := t.MAXPAGENO

	totalforce := t.Totalforce
	// 0 means no total calculation, this is the default i.e. no report of total number of pages
	// totalforce is not allowed to pass in args for securit reason
	// -1 means total number of pages is calculated from the database
	if totalforce == 0 || args[namePagesize] == nil || args[namePageno] != nil {
		return nil
	}

	nt := 0
	if totalforce < -1 { // take the absolute as the total number
		nt = int(math.Abs(float64(totalforce)))
	} else if totalforce == -1 || args[nameTotalno] == nil {
		if err := table.totalHashContext(ctx, db, &nt, extra...); err != nil {
			return err
		}
	} else {
		switch v := args[nameTotalno].(type) {
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

	args[nameTotalno] = nt
	nr := 0
	switch v := args[namePagesize].(type) {
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
	args[nameMaxpageno] = (nt-1)/nr + 1
	return nil
}

func (t *Topics) RunAction(db *sql.DB, table *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	return t.RunActionContext(context.Background(), db, table, args, extra...)
}

func (t *Topics) RunActionContext(ctx context.Context, db *sql.DB, table *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	t.setDefaultElementNames()
	sql, labels := table.filterPars(args, t.FIELDS, t.getAllowed())
	order := t.orderString(table, args)

	err := t.pagination(ctx, db, table, args, extra...)
	if err != nil {
		return nil, err
	}

	newExtra := table.byConstraint(args, extra...)
	if hasValue(newExtra) {
		where, values := selectCondition(newExtra, table.TableName)
		if where != "" {
			sql += "\nWHERE " + where
		}
		if order != "" {
			sql += "\n" + order
		}
		if table.dbDriver == Postgres {
			sql = questionMarkerNumber(sql)
		}

		return getSQL(ctx, db, table.logger, sql, labels, values...)
	}

	if order != "" {
		sql += "\n" + order
	}
	if table.dbDriver == Postgres {
		sql = questionMarkerNumber(sql)
	}

	return getSQL(ctx, db, table.logger, sql, labels)
}
