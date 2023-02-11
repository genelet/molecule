package protein

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/genelet/molecule/godbi"
)

const (
	PubDEFAULT   = "public"
	AdminDEFAULT = "admin"
	SignPOSTFIX  = "_sign"
)

type Protein struct {
	DBDriver     godbi.DBType           `json:"dbdriver,omitempty"`
	Teams        map[string]*Team       `json:"teams,omitempty"`
}

func (self *Protein) RunContext(ctx context.Context, db *sql.DB, token, team, atom, action string, args, extra map[string]interface{}) ([]interface{}, error) {
	if self.Teams == nil {
		return nil, fmt.Errorf("no teams found")
	}
	teamObj, ok := self.Teams[team]
	if !ok {
		return nil, fmt.Errorf("no team object found for %s", team)
	}

	return teamObj.RunContext(ctx, db, self.DBDriver, token, atom, action, args, extra)
}

func NewDefaultProtein(driver godbi.DBType, userIDs map[string]string) *Protein {
	teams := map[string]*Team{
		PubDEFAULT: &Team{IsPublic:true},
		AdminDEFAULT: &Team{IsAdmin:true}}
	for k, v := range userIDs {
		teams[k] = &Team{UserIDName: v}
	}
	
	return &Protein{driver, teams}
}
