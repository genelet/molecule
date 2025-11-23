package gometa

import (
	"github.com/genelet/molecule/godbi"
)

// MoleculeToGraph translates molecule into protobuf. args:
//
//   - oneofs format: map[atomName][oneofName][list of fields]
//   - package name
//   - goPackage name
//   - primary table name
//   - primary table's pk
//   - child to parent table mapping
//   - table name to pk mapping
func MoleculeToGraph(molecule *godbi.Molecule, args ...any) *Graph {
	var oneofs map[string]map[string][]string
	var packageName, goPackageName, pkTable, pkName string
	var pksTable, pks map[string]string
	if args != nil {
		if args[0] != nil {
			oneofs = args[0].(map[string]map[string][]string)
		}
		if len(args) >= 2 && args[1] != nil {
			packageName = args[1].(string)
		}
		if len(args) >= 3 && args[2] != nil {
			goPackageName = args[2].(string)
		}
		if len(args) >= 4 && args[3] != nil {
			pkTable = args[3].(string)
		}
		if len(args) >= 5 && args[4] != nil {
			pkName = args[4].(string)
		}
		if len(args) >= 6 && args[5] != nil {
			pksTable = args[5].(map[string]string)
		}
		if len(args) >= 7 && args[6] != nil {
			pks = args[6].(map[string]string)
		}
	}
	var nodes []*Node
	for _, atom := range molecule.Atoms {
		var node *Node
		if oneofs != nil {
			node = atomToNode(atom, oneofs[atom.AtomName])
		} else {
			node = atomToNode(atom)
		}
		nodes = append(nodes, node)
	}

	return &Graph{PackageName: packageName, PkTable: pkTable, PkName: pkName, GoPackageName: goPackageName, DBDriver: int32(molecule.DBDriver), PksTable: pksTable, Pks: pks, Nodes: nodes}
}

func atomToNode(atom *godbi.Atom, oneofs ...map[string][]string) *Node {
	nodeTable := atomTableToNodeTable(atom.Table, oneofs...)
	nodeActions := atomActionsToNodeActions(atom)
	return &Node{AtomName: atom.AtomName, AtomTable: nodeTable, AtomActions: nodeActions}
}

func getOneof(name string, oneofs ...map[string][]string) string {
	if oneofs == nil || oneofs[0] == nil {
		return ""
	}
	for k, vs := range oneofs[0] {
		for _, v := range vs {
			if v == name {
				return k
			}
		}
	}
	return ""
}

func atomTableToNodeTable(table godbi.Table, oneofs ...map[string][]string) *Node_Table {
	nodeTable := &Node_Table{}

	nodeTable.TableName = table.TableName
	for _, col := range table.Columns {
		nodeCol := &Node_Table_Col{
			ColumnName: col.ColumnName,
			TypeName:   col.TypeName,
			Label:      col.Label,
			InOneof:    getOneof(col.ColumnName, oneofs...),
			Notnull:    col.Notnull,
			Constraint: col.Constraint,
			Auto:       col.Auto,
			Recurse:    col.Recurse}
		nodeTable.Columns = append(nodeTable.Columns, nodeCol)
	}
	nodeTable.Pks = table.Pks
	nodeTable.IDAuto = table.IDAuto
	for _, fk := range table.Fks {
		nodeFk := &Node_Table_Fk{
			FkTable:  fk.FkTable,
			FkColumn: fk.FkColumn,
			Column:   fk.Column}
		nodeTable.Fks = append(nodeTable.Fks, nodeFk)
	}
	nodeTable.Uniques = table.Uniques

	return nodeTable
}

func atomActionsToNodeActions(atom *godbi.Atom) *Node_Actions {
	nodeActions := &Node_Actions{}

	dbiConnection := func(conn *godbi.Connection) *Node_Actions_Connection {
		return &Node_Actions_Connection{
			AtomName:    conn.AtomName,
			ActionName:  conn.ActionName,
			Dimension:   Node_Actions_ConnectType(conn.Dimension),
			Marker:      conn.Marker,
			RelateArgs:  conn.RelateArgs,
			RelateExtra: conn.RelateExtra}
	}

	if insert := atom.GetAction("insert"); insert != nil {
		nodeInsert := &Node_Actions_Insert{
			ActionName: "insert",
			Picked:     insert.(*godbi.Insert).Picked,
			IsDo:       insert.GetBaseAction().IsDo}
		for _, prepare := range insert.GetBaseAction().Prepares {
			nodeInsert.PrepareConnects = append(nodeInsert.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range insert.GetBaseAction().Nextpages {
			nodeInsert.NextpageConnects = append(nodeInsert.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.InsertItem = nodeInsert
	}

	if insupd := atom.GetAction("insupd"); insupd != nil {
		nodeInsupd := &Node_Actions_Insupd{
			ActionName: "insupd",
			Picked:     insupd.(*godbi.Insupd).Picked,
			IsDo:       insupd.GetBaseAction().IsDo}
		for _, prepare := range insupd.GetBaseAction().Prepares {
			nodeInsupd.PrepareConnects = append(nodeInsupd.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range insupd.GetBaseAction().Nextpages {
			nodeInsupd.NextpageConnects = append(nodeInsupd.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.InsupdItem = nodeInsupd
	}

	if iupdate := atom.GetAction("update"); iupdate != nil {
		update := iupdate.(*godbi.Update)
		nodeUpdate := &Node_Actions_Update{
			ActionName: "update",
			IsDo:       update.IsDo,
			Picked:     update.Picked,
			Empties:    update.Empties}
		for _, prepare := range update.Prepares {
			nodeUpdate.PrepareConnects = append(nodeUpdate.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range update.Nextpages {
			nodeUpdate.NextpageConnects = append(nodeUpdate.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.UpdateItem = nodeUpdate
	}

	if delett := atom.GetAction("delete"); delett != nil {
		nodeDelete := &Node_Actions_Delete{ActionName: "delete", IsDo: delett.GetBaseAction().IsDo}
		for _, prepare := range delett.GetBaseAction().Prepares {
			nodeDelete.PrepareConnects = append(nodeDelete.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range delett.GetBaseAction().Nextpages {
			nodeDelete.NextpageConnects = append(nodeDelete.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.DeleteItem = nodeDelete
	}

	if delecs := atom.GetAction("delecs"); delecs != nil {
		nodeDelecs := &Node_Actions_Delecs{ActionName: "delecs", IsDo: delecs.GetBaseAction().IsDo}
		for _, prepare := range delecs.GetBaseAction().Prepares {
			nodeDelecs.PrepareConnects = append(nodeDelecs.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range delecs.GetBaseAction().Nextpages {
			nodeDelecs.NextpageConnects = append(nodeDelecs.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.DelecsItem = nodeDelecs
	}

	if itopics := atom.GetAction("topics"); itopics != nil {
		topics := itopics.(*godbi.Topics)
		nodeTopics := &Node_Actions_Topics{
			IsDo:        topics.IsDo,
			ActionName:  "topics",
			Picked:      topics.Picked,
			FIELDS:      topics.FIELDS,
			Totalforce:  int32(topics.Totalforce),
			MAXPAGENO:   topics.MAXPAGENO,
			TOTALNO:     topics.TOTALNO,
			PAGESIZE:    topics.PAGESIZE,
			PAGENO:      topics.PAGENO,
			SORTBY:      topics.SORTBY,
			SORTREVERSE: topics.SORTREVERSE}
		for _, prepare := range topics.Prepares {
			nodeTopics.PrepareConnects = append(nodeTopics.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range topics.Nextpages {
			nodeTopics.NextpageConnects = append(nodeTopics.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.TopicsItem = nodeTopics
	}

	if iedit := atom.GetAction("edit"); iedit != nil {
		edit := iedit.(*godbi.Edit)
		nodeEdit := &Node_Actions_Edit{
			IsDo:       edit.IsDo,
			ActionName: "edit",
			Picked:     edit.Picked,
			FIELDS:     edit.FIELDS}
		for _, prepare := range edit.Prepares {
			nodeEdit.PrepareConnects = append(nodeEdit.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range edit.Nextpages {
			nodeEdit.NextpageConnects = append(nodeEdit.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.EditItem = nodeEdit
	}

	return nodeActions
}
