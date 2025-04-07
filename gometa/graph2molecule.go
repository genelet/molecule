package gometa

import (
	"github.com/genelet/molecule/godbi"
)

// GraphToMolecule translates protobuf message into molecule with possible oneof map
//
// oneof format: map[atomName][oneofName][list of fields]
func GraphToMolecule(graph *Graph) (*godbi.Molecule, map[string]map[string][]string) {
	var atoms []*godbi.Atom
	var oneofs map[string]map[string][]string
	for _, node := range graph.Nodes {
		atom, hash := nodeToAtom(node)
		atoms = append(atoms, atom)
		if hash != nil {
			if oneofs == nil {
				oneofs = make(map[string]map[string][]string)
			}
			oneofs[node.AtomTable.TableName] = hash
		}
	}
	return &godbi.Molecule{Atoms: atoms, DBDriver: godbi.DBType(graph.DBDriver)}, oneofs
}

func nodeToAtom(node *Node) (*godbi.Atom, map[string][]string) {
	atomTable, oneofs := nodeTableToAtomTable(node.AtomTable)
	atomActions := nodeActionsToAtomActions(node.AtomActions)
	return &godbi.Atom{AtomName: node.AtomName, Table: *atomTable, Actions: atomActions}, oneofs
}

func nodeTableToAtomTable(nodeTable *Node_Table) (*godbi.Table, map[string][]string) {
	atomTable := &godbi.Table{}
	atomTable.TableName = nodeTable.GetTableName()
	atomTable.Pks = nodeTable.GetPks()
	atomTable.IDAuto = nodeTable.GetIDAuto()
	atomTable.Uniques = nodeTable.GetUniques()

	var oneofs map[string][]string

	for _, col := range nodeTable.GetColumns() {
		atomCol := &godbi.Col{
			ColumnName: col.GetColumnName(),
			TypeName:   col.GetTypeName(),
			Label:      col.GetLabel(),
			Notnull:    col.GetNotnull(),
			Constraint: col.GetConstraint(),
			Auto:       col.GetAuto(),
			Recurse:    col.GetRecurse()}
		atomTable.Columns = append(atomTable.Columns, atomCol)
		group := col.GetInOneof()
		if group != "" {
			if oneofs == nil {
				oneofs = make(map[string][]string)
			}
			if oneofs[group] == nil {
				oneofs[group] = make([]string, 0)
			}
			oneofs[group] = append(oneofs[group], col.GetColumnName())
		}
	}
	for _, fk := range nodeTable.GetFks() {
		atomFk := &godbi.Fk{
			FkTable:  fk.GetFkTable(),
			FkColumn: fk.GetFkColumn(),
			Column:   fk.GetColumn()}
		atomTable.Fks = append(atomTable.Fks, atomFk)
	}

	return atomTable, oneofs
}

func nodeActionsToAtomActions(nodeActions *Node_Actions) []godbi.Capability {
	var actions []godbi.Capability

	dbiConnection := func(conn *Node_Actions_Connection) *godbi.Connection {
		return &godbi.Connection{
			AtomName:    conn.GetAtomName(),
			ActionName:  conn.GetActionName(),
			Dimension:   godbi.ConnectType(conn.GetDimension()),
			Marker:      conn.GetMarker(),
			RelateArgs:  conn.GetRelateArgs(),
			RelateExtra: conn.GetRelateExtra()}
	}

	if insert := nodeActions.GetInsertItem(); insert != nil {
		atomInsert := &godbi.Insert{Action: godbi.Action{ActionName: "insert"}}
		atomInsert.SetIsDo(insert.GetIsDo())
		atomInsert.Picked = insert.GetPicked()
		for _, prepare := range insert.GetPrepareConnects() {
			atomInsert.Prepares = append(atomInsert.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range insert.GetNextpageConnects() {
			atomInsert.Nextpages = append(atomInsert.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomInsert)
	}

	if insupd := nodeActions.GetInsupdItem(); insupd != nil {
		atomInsupd := &godbi.Insupd{Action: godbi.Action{ActionName: "insupd"}}
		atomInsupd.SetIsDo(insupd.GetIsDo())
		atomInsupd.Picked = insupd.GetPicked()
		for _, prepare := range insupd.GetPrepareConnects() {
			atomInsupd.Prepares = append(atomInsupd.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range insupd.GetNextpageConnects() {
			atomInsupd.Nextpages = append(atomInsupd.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomInsupd)
	}

	if update := nodeActions.GetUpdateItem(); update != nil {
		atomUpdate := &godbi.Update{Action: godbi.Action{ActionName: "update"}}
		atomUpdate.SetIsDo(update.GetIsDo())
		atomUpdate.Empties = update.GetEmpties()
		atomUpdate.Picked = update.GetPicked()
		for _, prepare := range update.GetPrepareConnects() {
			atomUpdate.Prepares = append(atomUpdate.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range update.GetNextpageConnects() {
			atomUpdate.Nextpages = append(atomUpdate.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomUpdate)
	}

	if delett := nodeActions.GetDeleteItem(); delett != nil {
		atomDelete := &godbi.Delete{Action: godbi.Action{ActionName: "delete"}}
		atomDelete.SetIsDo(delett.GetIsDo())
		for _, prepare := range delett.GetPrepareConnects() {
			atomDelete.Prepares = append(atomDelete.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range delett.GetNextpageConnects() {
			atomDelete.Nextpages = append(atomDelete.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomDelete)
	}

	if delecs := nodeActions.GetDelecsItem(); delecs != nil {
		atomDelecs := &godbi.Delecs{Action: godbi.Action{ActionName: "delecs"}}
		atomDelecs.SetIsDo(delecs.GetIsDo())
		for _, prepare := range delecs.GetPrepareConnects() {
			atomDelecs.Prepares = append(atomDelecs.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range delecs.GetNextpageConnects() {
			atomDelecs.Nextpages = append(atomDelecs.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomDelecs)
	}

	if topics := nodeActions.GetTopicsItem(); topics != nil {
		atomTopics := &godbi.Topics{Action: godbi.Action{ActionName: "topics"}}
		atomTopics.SetIsDo(topics.GetIsDo())
		atomTopics.FIELDS = topics.GetFIELDS()
		atomTopics.Totalforce = int(topics.GetTotalforce())
		atomTopics.MAXPAGENO = topics.GetMAXPAGENO()
		atomTopics.TOTALNO = topics.GetTOTALNO()
		atomTopics.PAGESIZE = topics.GetPAGESIZE()
		atomTopics.PAGENO = topics.GetPAGENO()
		atomTopics.SORTBY = topics.GetSORTBY()
		atomTopics.SORTREVERSE = topics.GetSORTREVERSE()
		atomTopics.Picked = topics.GetPicked()
		for _, prepare := range topics.GetPrepareConnects() {
			atomTopics.Prepares = append(atomTopics.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range topics.GetNextpageConnects() {
			atomTopics.Nextpages = append(atomTopics.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomTopics)
	}

	if edit := nodeActions.GetEditItem(); edit != nil {
		atomEdit := &godbi.Edit{Action: godbi.Action{ActionName: "edit"}}
		atomEdit.SetIsDo(atomEdit.GetIsDo())
		atomEdit.FIELDS = edit.GetFIELDS()
		atomEdit.Picked = edit.GetPicked()
		for _, prepare := range edit.GetPrepareConnects() {
			atomEdit.Prepares = append(atomEdit.Prepares, dbiConnection(prepare))
		}
		for _, nextpage := range edit.GetNextpageConnects() {
			atomEdit.Nextpages = append(atomEdit.Nextpages, dbiConnection(nextpage))
		}
		actions = append(actions, atomEdit)
	}

	return actions
}
