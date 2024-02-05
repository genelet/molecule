package godbi

import (
	"context"
	"database/sql"
)

const (
	STMT = "stmt"
)

type StmtContext struct {
	Statement string        `json:"statement" hcl:"statement,optional"`
	Labels    []interface{} `json:"labels,omitempty" hcl:"labels,optional"`
	Pars      []string      `json:"pars,omitempty" hcl:"pars,optional"`
}

// Stmt struct for search one specific row by primary key
type Stmt struct {
	Action
	StmtContext
	Others map[string]*StmtContext `json:"others,omitempty" hcl:"others,block"`
}

var _ Capability = (*Stmt)(nil)

func (self *Stmt) RunActionContext(ctx context.Context, db *sql.DB, t *Table, ARGS map[string]interface{}, extra ...map[string]interface{}) ([]interface{}, error) {
	var statement string
	var pars []string
	var labels []interface{}
	if v, ok := ARGS[STMT]; ok {
		if self.Others == nil || self.Others[v.(string)] == nil {
			return nil, errorActionNil(v.(string))
		}
		other := self.Others[v.(string)]
		statement = other.Statement
		pars = other.Pars
		labels = other.Labels
	} else {
		statement = self.Statement
		pars = self.Pars
		labels = self.Labels
		if labels == nil {
			for _, col := range self.Picked {
				labels = append(labels, col)
			}
		}
	}

	var extra0 map[string]interface{}
	if extra != nil {
		extra0 = extra[0]
	}
	ids := properValues(pars, ARGS, extra0)

	return getSQL(ctx, db, statement, labels, ids...)
}
