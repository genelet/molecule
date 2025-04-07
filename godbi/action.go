package godbi

import (
	"context"
	"database/sql"
)

// Capability is to implement Capability interface
type Capability interface {
	GetActionName() string
	SetActionName(string)
	GetPrepares() []*Connection
	GetNextpages() []*Connection
	GetIsDo(...string) bool
	SetIsDo(bool)
	SetPrepares([]*Connection)
	SetNextpages([]*Connection)
	GetPicked() []string
	SetPicked([]string)
	// RunActionContext runs the action with context, db, table, and args
	RunActionContext(context.Context, *sql.DB, *Table, map[string]any, ...map[string]any) ([]any, error)
}

// Action is the base struct for REST actions. Prepares and Nextpages are edges to other tables before and after the action.
type Action struct {
	ActionName string        `json:"actionName,omitempty" hcl:"actionName,label"`
	Picked     []string      `json:"picked,omitempty" hcl:"picked,optional"`
	Prepares   []*Connection `json:"prepares,omitempty" hcl:"prepares,block"`
	Nextpages  []*Connection `json:"nextpages,omitempty" hcl:"nextpages,block"`
	isDo       bool
}

// GetActionName gets the action name
func (self *Action) GetActionName() string {
	return self.ActionName
}

// GetPrepares gets the prepares
func (self *Action) GetPrepares() []*Connection {
	return self.Prepares
}

// GetNextpages gets the nextpages
func (self *Action) GetNextpages() []*Connection {
	return self.Nextpages
}

// GetIsDo gets the isDo
func (self *Action) GetIsDo(option ...string) bool {
	if option != nil {
		switch option[0] {
		case "insert", "insupd", "update", "delete", "delecs":
			return true
		default:
		}
	}
	return self.isDo
}

// SetActionName sets the action name
func (self *Action) SetActionName(name string) {
	self.ActionName = name
}

// SetIsDo sets the isDo
func (self *Action) SetIsDo(is bool) {
	self.isDo = is
}

// GetPicked gets the picked
func (self *Action) GetPicked() []string {
	return self.Picked
}

// SetPicked sets the picked
func (self *Action) SetPicked(picked []string) {
	self.Picked = picked
}

func (self *Action) getAllowed() map[string]bool {
	if self.Picked == nil {
		return nil
	}
	allowed := make(map[string]bool)
	for _, v := range self.Picked {
		allowed[v] = true
	}
	return allowed
}

// SetPrepares sets the prepares
func (self *Action) SetPrepares(x []*Connection) {
	self.Prepares = x
}

// SetNextpages sets the nextpages
func (self *Action) SetNextpages(x []*Connection) {
	self.Nextpages = x
}

func getSQL(ctx context.Context, db *sql.DB, logger Slogger, statement string, labels []any, ids ...any) ([]any, error) {
	lists := make([]any, 0)
	dbi := &DBI{DB: db, logger: logger}
	err := dbi.SelectSQLContext(ctx, &lists, statement, labels, ids...)
	return lists, err
}
