package godbi

import (
	"context"
	"database/sql"
)

const (
	STMT = "stmt"
)

type StmtContext struct {
	Statement string   `json:"statement" hcl:"statement,optional"`
	Labels    []any    `json:"labels,omitempty" hcl:"labels,optional"`
	Pars      []string `json:"pars,omitempty" hcl:"pars,optional"`
}

// Stmt struct for search one specific row by primary key
type Stmt struct {
	Action
	StmtContext
	Others map[string]*StmtContext `json:"others,omitempty" hcl:"others,block"`
}

var _ Capability = (*Stmt)(nil)

func (s *Stmt) RunActionContext(ctx context.Context, db *sql.DB, t *Table, args map[string]any, extra ...map[string]any) ([]any, error) {
	var statement string
	var pars []string
	var labels []any
	if v, ok := args[STMT]; ok {
		if s.Others == nil || s.Others[v.(string)] == nil {
			return nil, errorActionNil(v.(string))
		}
		other := s.Others[v.(string)]
		statement = other.Statement
		pars = other.Pars
		labels = other.Labels
	} else {
		statement = s.Statement
		pars = s.Pars
		labels = s.Labels
		if labels == nil {
			for _, col := range s.Picked {
				labels = append(labels, col)
			}
		}
	}

	var extra0 map[string]any
	if extra != nil {
		extra0 = extra[0]
	}
	ids := properValues(pars, args, extra0)

	return getSQL(ctx, db, t.logger, statement, labels, ids...)
}
