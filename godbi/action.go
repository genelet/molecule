package godbi

import (
	"context"
	"database/sql"
)

// Action is to implement Capability interface
//
type Capability interface {
	GetActionName() string
	GetPrepares()  []*Connection
	GetNextpages() []*Connection
	GetIsDo() bool
	GetAppendix() interface{}
	SetPrepares([]*Connection)
	SetNextpages([]*Connection)
	SetAppendix(interface{})
	RunActionContext(context.Context, *sql.DB, *Table, map[string]interface{}, ...map[string]interface{}) ([]interface{}, error)
}

type Action struct {
	ActionName string `json:"actionName,omitempty" hcl:"actionName,optional"`
	Prepares  []*Connection `json:"Prepares,omitempty" hcl:"Prepares,block"`
	Nextpages []*Connection `json:"nextpages,omitempty" hcl:"nextpages,block"`
	IsDo      bool          `json:"isDo,omitempty" hcl:"isDo,optional"`
	Appendix  interface{}   `json:"appendix,omitempty" hcl:"appendix,block"`
}

func (self *Action) GetActionName() string {
	return self.ActionName
}

func (self *Action) GetPrepares() []*Connection {
	return self.Prepares
}

func (self *Action) GetNextpages() []*Connection {
	return self.Nextpages
}

func (self *Action) GetIsDo() bool {
	return self.IsDo
}

func (self *Action) GetAppendix() interface{} {
	return self.Appendix
}

func (self *Action) SetPrepares(x []*Connection) {
	self.Prepares = x
}

func (self *Action) SetNextpages(x []*Connection) {
	self.Nextpages = x
}

func (self *Action) SetAppendix(x interface{}) {
	self.Appendix = x
}
