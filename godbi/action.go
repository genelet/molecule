package godbi

import (
	"context"
	"database/sql"
)

// Capability is to implement Capability interface
type Capability interface {
	GetBaseAction() *Action
	// RunActionContext runs the action with context, db, table, and args
	RunActionContext(context.Context, *sql.DB, *Table, map[string]any, ...map[string]any) ([]any, error)
}

// Action is the base struct for REST actions. Prepares and Nextpages are edges to other tables before and after the action.
type Action struct {
	ActionName string        `json:"actionName,omitempty" hcl:"actionName,label"`
	Picked     []string      `json:"picked,omitempty" hcl:"picked,optional"`
	Prepares   []*Connection `json:"prepares,omitempty" hcl:"prepares,block"`
	Nextpages  []*Connection `json:"nextpages,omitempty" hcl:"nextpages,block"`
	IsDo       bool          `json:"-" hcl:"-"`
}

// GetBaseAction gets the base action
func (a *Action) GetBaseAction() *Action {
	return a
}

func (a *Action) getAllowed() map[string]bool {
	if a.Picked == nil {
		return nil
	}
	allowed := make(map[string]bool)
	for _, v := range a.Picked {
		allowed[v] = true
	}
	return allowed
}

func getSQL(ctx context.Context, db *sql.DB, logger Slogger, statement string, labels []any, ids ...any) ([]any, error) {
	lists := make([]any, 0)
	dbi := &DBI{DB: db, logger: logger}
	err := dbi.SelectSQLContext(ctx, &lists, statement, labels, ids...)
	return lists, err
}
