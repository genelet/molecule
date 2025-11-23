package rdb

import (
	"database/sql"

	"github.com/genelet/molecule/godbi"
)

type schema interface {
	tableNames(*sql.DB) ([]string, error)
	getTable(*sql.DB, string) (*godbi.Table, error)
	getFks(*sql.DB, string) ([]*godbi.Fk, error)
}

type database struct {
	schema
	DBDriver     godbi.DBType
	DatabaseName string
}

func (d *database) GetMolecule(db *sql.DB) (*godbi.Molecule, error) {
	scheme := d.schema

	tableNames, err := scheme.tableNames(db)
	if err != nil {
		return nil, err
	}

	refPk := make(map[string]string)

	var atoms []*godbi.Atom
	nextpages := make(map[string]map[string][]*godbi.Connection)
	prepares := make(map[string]map[string][]*godbi.Connection)
	for _, name := range tableNames {
		table, err := scheme.getTable(db, name)
		if err != nil {
			return nil, err
		}
		if table.Pks != nil && len(table.Pks) == 1 {
			refPk[table.TableName] = table.Pks[0]
		}
		atom := autoAtom(table)
		fks, err := scheme.getFks(db, name)
		if err != nil {
			return nil, err
		}
		atom.Fks = fks
		for _, fk := range fks {
			// this is one of possible choices
			// pls check adv_balance in adv_campaign, when viewed as "team"
			autoConnection(table, fk, nextpages, prepares)
		}
		atoms = append(atoms, atom)
	}

	var newAtoms []*godbi.Atom
	for _, atom := range atoms {
		actions := setConnections(atom, nextpages, prepares)
		newAtoms = append(newAtoms, &godbi.Atom{AtomName: atom.AtomName, Table: atom.Table, Actions: actions})
	}

	return &godbi.Molecule{Atoms: newAtoms, DBDriver: d.DBDriver}, nil
}

func autoAtom(table *godbi.Table) *godbi.Atom {
	edit := new(godbi.Edit)
	edit.ActionName = "edit"
	topics := new(godbi.Topics)
	topics.ActionName = "topics"
	insert := new(godbi.Insert)
	insert.ActionName = "insert"
	insert.IsDo = true
	update := new(godbi.Update)
	update.ActionName = "update"
	update.IsDo = true
	insupd := new(godbi.Insupd)
	insupd.ActionName = "insupd"
	insupd.IsDo = true
	delett := new(godbi.Delete)
	delett.ActionName = "delete"
	delett.IsDo = true
	capas := []godbi.Capability{edit, topics, insert, update, insupd, delett}
	if table.IDAuto != "" {
		delecs := new(godbi.Delecs)
		delecs.ActionName = "delecs"
		delecs.IsDo = true
		delecs.Nextpages = []*godbi.Connection{{
			AtomName:   table.TableName,
			ActionName: "delete",
			RelateArgs: map[string]string{table.IDAuto: table.IDAuto}}}
		capas = append(capas, delecs)
	}
	return &godbi.Atom{AtomName: table.TableName, Table: *table, Actions: capas}
}

func autoConnection(table *godbi.Table, fk *godbi.Fk, nextpages, prepares map[string]map[string][]*godbi.Connection) {
	patom := fk.FkTable
	tatom := table.TableName
	if nextpages[patom] == nil {
		nextpages[patom] = make(map[string][]*godbi.Connection)
	}
	if prepares[patom] == nil {
		prepares[patom] = make(map[string][]*godbi.Connection)
	}
	if prepares[tatom] == nil {
		prepares[tatom] = make(map[string][]*godbi.Connection)
	}
	for _, actionName := range []string{"topics", "edit"} {
		def := "topics"
		nextpage := &godbi.Connection{AtomName: tatom, ActionName: def, Marker: tatom}
		nextpage.RelateExtra = map[string]string{fk.FkColumn: fk.Column}
		nextpages[patom][actionName] = append(nextpages[patom][actionName], nextpage)
	}
	for _, actionName := range []string{"insert", "insupd", "update"} {
		nextpage := &godbi.Connection{AtomName: tatom, ActionName: actionName, Marker: tatom}
		nextpage.RelateArgs = map[string]string{fk.FkColumn: fk.Column}
		nextpages[patom][actionName] = append(nextpages[patom][actionName], nextpage)

		prepare := &godbi.Connection{AtomName: patom, ActionName: actionName, Marker: patom}
		prepare.RelateArgs = map[string]string{fk.FkColumn: fk.Column}
		prepares[tatom][actionName] = append(prepares[tatom][actionName], prepare)
	}

	prepare := &godbi.Connection{AtomName: tatom, ActionName: "delecs", RelateArgs: map[string]string{fk.FkColumn: fk.Column}}
	prepares[patom]["delete"] = append(prepares[patom]["delete"], prepare)
}

func setConnections(atom *godbi.Atom, nextpages, prepares map[string]map[string][]*godbi.Connection) []godbi.Capability {
	atomName := atom.Table.TableName
	actions := atom.Actions
	for tableName, actionMap := range nextpages {
		if tableName != atomName {
			continue
		}
		for actionName, nextpages := range actionMap {
			for _, action := range actions {
				if actionName == action.GetBaseAction().ActionName {
					action.GetBaseAction().Nextpages = nextpages
				}
			}
		}
	}
	for tableName, actionMap := range prepares {
		if tableName != atomName {
			continue
		}
		for actionName, prepares := range actionMap {
			for _, action := range actions {
				if actionName == action.GetBaseAction().ActionName {
					action.GetBaseAction().Prepares = prepares
				}
			}
		}
	}
	return actions
}

// toString helper function to safely convert interface{} to string
func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
