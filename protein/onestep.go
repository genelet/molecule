package protein

import (
	"github.com/genelet/molecule/godbi"
)

type OneStepStopper struct {
	component string
	token     string
}

func (self *OneStepStopper) Sign(tableObj *godbi.Table, item interface{}) bool {
	if self.component != tableObj.TableName {
		return false
	}

	if item == nil {
		return true
	}
	hash, ok := item.(map[string]interface{})
	if !ok {
		return false
	}

	if self.token != "" {
		pk := tableObj.Pks[0]
		if val, ok := hash[pk]; ok {
			hash[pk+SIGN] = digest(self.token, tableObj.TableName, pk, val)
		}
	}
	return true
}
