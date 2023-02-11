package protein

import (
	"fmt"
	"github.com/genelet/molecule/godbi"
)

// AutoTeams assigns colorfuls to teams using atoms and foreign keys 
// in molecule. The teams have to exist already.
func (self *Protein) AutoTeams(molecule *godbi.Molecule, userTables map[string]string) error {
	if self.Teams == nil {
		return fmt.Errorf("teams not defined")
	}
	ref := make(map[string]bool)
	var pubName, adminName string

	for teamName, team := range self.Teams {
		if team.IsPublic {
			pubName = teamName
			continue
		} else if team.IsAdmin {
			adminName = teamName
			team.AutoUserTable(molecule)
		} else {
			team.AutoUserTable(molecule, userTables[teamName])
			for _, colorful := range team.Colorfuls {
				tableName := colorful.Atom.TableName
				if colorful.IsProtect() {
					ref[tableName] = true
				}
			}
		}
		self.Teams[teamName] = team
	}

	if pubName == "" || adminName == "" {
		return fmt.Errorf("public or admin team not defined")
	}

	var colorfuls []*Colorful
	for _, colorful := range self.Teams[adminName].Colorfuls {
		tableName := colorful.Atom.TableName
		if _, ok := ref[tableName]; ok { continue }
		colorfuls = append(colorfuls, &Colorful{Atom: colorful.Atom})
	}
	self.Teams[pubName].Colorfuls = colorfuls

	return nil
}

// AutoUserTable assigns colorfuls to team. In case of admin, all atoms 
// are becoming colorfuls.
//  - userTable
//  - bool for backward
func (self *Team) AutoUserTable(molecule *godbi.Molecule, args ...interface{}) {
	if self.IsPublic {
		return
	}
	var colorfuls []*Colorful
	for _, atom := range molecule.Atoms {
		colorfuls = append(colorfuls, &Colorful{Atom: atom.(*godbi.Atom)})
	}
	if self.IsAdmin {
		self.Colorfuls = colorfuls
		return
	}

	if len(args) < 1 { return }
	userTable := args[0].(string)
	withBackward := false
	if len(args) > 1 && args[1].(bool) == true { withBackward = true }

	colors := userTableProps(molecule, userTable, self.UserIDName, withBackward)
	if colors == nil || len(colors) == 0 { return }

	for _, colorful := range colorfuls {
		if color, ok := colors[colorful.Atom.GetTable().TableName]; ok {
			colorful.Color = *color
		}
		self.Colorfuls = append(self.Colorfuls, colorful)
	}
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

func userTableProps(molecule *godbi.Molecule, userTable, userid string, withBackward bool) map[string]*Color {
	forwards := forwardReference(molecule)

	colors := make(map[string]*Color)
	ref := make(map[string][]*godbi.Fk)

	for _, atom := range molecule.Atoms {
		table := atom.GetTable()
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

	for {
		var step1, step2 int
		// forward keys (in ads)
		for tableName, hash := range colors {
			if hash.IsUser || !hash.IsLanding {
				continue
			}
			step1 += assignProps(tableName, colors, ref)
		}
		if withBackward == false { break }

		// backward keys (e.g. in proto). loop back to forward
		for tableName, fks := range ref {
			_, ok := colors[tableName]
			if ok { continue }
			for _, fk := range fks {
				_, ok := colors[fk.FkTable]
				if ok {
					colors[tableName] = NewColor(false, false, fk.Column)
					step2 += 1
					break
				}
			}
		}
		if step1 == 0 && step2 == 0 { break }
	}

	return colors
}

func assignProps(tableName string, colors map[string]*Color, ref map[string][]*godbi.Fk) int {
	offsprings, ok := ref[tableName]
	if !ok {
		return 0
	}
	num := 0
	for _, offspring := range offsprings {
		tname := offspring.FkTable
		_, ok := colors[tname]
		if ok { continue }
		colors[tname] = NewColor(false, false, offspring.FkColumn)
		num++
		num += assignProps(tname, colors, ref)
	}
	return num
}
