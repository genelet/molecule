package protein

import (
	"github.com/genelet/molecule/godbi"
)

func NewTeamUserTable(molecule *godbi.Molecule, teamName string, users ...string) *Team {
	var userTable string
	var team *Team
	if users == nil {
		team = &Team{IsAdmin: true}
	} else {
		team = &Team{UserIDName: users[0]}
		if len(users) == 2 {
			userTable = users[1]
		} else {
			userTable = teamName
		}
	}
	team.AutoUserTable(molecule, teamName, userTable)

	return team
}

func (self *Team) AutoUserTable(molecule *godbi.Molecule, teamName, userTable string) {
	var colors map[string]*Color

	if !self.IsAdmin && !self.IsPublic {
		colors = userTableProps(molecule, userTable, self.UserIDName)
	}

	var colorfuls []*Colorful

	for _, atom := range molecule.Atoms {
		tableName := atom.GetTable().TableName
		colorful := &Colorful{Atom: atom.(*godbi.Atom)}
		if colors != nil {
			if color, ok := colors[tableName]; ok {
				colorful.Color = *color
			}
		}
		colorfuls = append(colorfuls, colorful)
	}

	self.Colorfuls = colorfuls
}

// for each table in molecule, this function lists its forwarded fks
// i.e. its fk.Column appears in fk.FkTable as foreign key fk.FkColumn
func forwardReference(molecule *godbi.Molecule) map[string][]*godbi.Fk {
	forward := make(map[string][]*godbi.Fk)
	for _, atom := range molecule.Atoms {
		table := atom.(*godbi.Atom).Table
		for _, fk := range table.Fks {
			fw := &godbi.Fk{FkTable: table.TableName, FkColumn: fk.Column, Column: fk.FkColumn}
			forward[fk.FkTable] = append(forward[fk.FkTable], fw)
		}
	}
	return forward
}

func userTableProps(molecule *godbi.Molecule, userTable, userid string) map[string]*Color {
	forwards := forwardReference(molecule)

	colors := make(map[string]*Color)
	ref := make(map[string][]*godbi.Fk)

	for _, atom := range molecule.Atoms {
		table := atom.(*godbi.Atom).Table
		tableName := table.TableName
		if table.Pks != nil && len(table.Pks) == 1 {
			pk := table.Pks[0]
			if tableName == userTable {
				if pk != userid {
					return nil
				}
				colors[tableName] = NewColor(true, false, userid)
				for _, fk := range forwards[tableName] {
					if pk == fk.Column {
						colors[fk.FkTable] = NewColor(false, true, fk.FkColumn)
					}
				}
			} else {
				for _, fk := range forwards[tableName] {
					if pk == fk.Column {
						ref[tableName] = append(ref[tableName], fk)
					}
				}
			}
		}
	}

	for tableName, hash := range colors {
		if hash.IsUser || !hash.IsLanding {
			continue
		}
		assignProps(tableName, colors, ref)
	}

	return colors
}

func assignProps(tableName string, colors map[string]*Color, ref map[string][]*godbi.Fk) {
	offsprings, ok := ref[tableName]
	if !ok {
		return
	}
	for _, offspring := range offsprings {
		tname := offspring.FkTable
		colors[tname] = NewColor(false, false, offspring.FkColumn)
		assignProps(tname, colors, ref)
	}
}
