package gometa

import (
	"github.com/genelet/molecule/godbi"
)

// MoleculeToGraph translates molecule (plus possible oneof map, and package name) into protobuf
//
// oneofs format: map[atomName][oneofName][list of fields]
//
func MoleculeToGraph(molecule *godbi.Molecule, rest ...interface{}) *Graph {
	var oneofs map[string]map[string][]string
	var packageName string
	if rest != nil && rest[0] != nil {
		oneofs = rest[0].(map[string]map[string][]string)
	}
	if rest != nil && rest[1] != nil {
		packageName = rest[1].(string)
	}
	var nodes []*Node
	for _, atom := range molecule.Atoms {
		var node *Node
		if oneofs != nil {
			node = atomToNode(atom, oneofs[atom.GetTable().TableName])
		} else {
			node = atomToNode(atom)
		}
		nodes = append(nodes, node)
	}

	return &Graph{PackageName: packageName, DatabaseName: molecule.DatabaseName, Nodes: nodes}
}

func atomToNode(atom godbi.Navigate, oneofs ...map[string][]string) *Node {
	nodeTable := atomTableToNodeTable(atom.GetTable(), oneofs...)
	nodeActions := atomActionsToNodeActions(atom)
	return &Node{AtomTable:nodeTable, AtomActions:nodeActions}
}

func getOneof(name string, oneofs ...map[string][]string) string {
	if oneofs == nil || oneofs[0] == nil { return "" }
	for k, vs := range oneofs[0] {
		for _, v := range vs {
			if v == name {
				return k
			}
		}
	}
	return ""
}

func atomTableToNodeTable(table *godbi.Table, oneofs ...map[string][]string) *Node_Table {
    nodeTable := &Node_Table{}

	nodeTable.TableName = table.TableName
	for _, col := range table.Columns {
		nodeCol := &Node_Table_Col{
			ColumnName: col.ColumnName,
			TypeName: col.TypeName,
			Label: col.Label,
			InOneof: getOneof(col.ColumnName, oneofs...),
			Notnull: col.Notnull,
			Auto: col.Auto,
			Recurse: col.Recurse}
		nodeTable.Columns = append(nodeTable.Columns, nodeCol)
	}
    nodeTable.Pks = table.Pks
	nodeTable.IdAuto = table.IdAuto
	for _, fk := range table.Fks {
		nodeFk := &Node_Table_Fk{
			FkTable: fk.FkTable,
			FkColumn: fk.FkColumn,
			Column: fk.Column}
		nodeTable.Fks = append(nodeTable.Fks, nodeFk)
	}
	nodeTable.Uniques = table.Uniques

	return nodeTable
}

func atomActionsToNodeActions(atom godbi.Navigate) *Node_Actions {
	nodeActions := &Node_Actions{}

	dbiConnection := func(conn *godbi.Connection) *Node_Actions_Connection {
		return &Node_Actions_Connection{
			TableName: conn.TableName,
			ActionName: conn.ActionName,
			Dimension: Node_Actions_ConnectType(conn.Dimension),
			Marker: conn.Marker,
			RelateArgs: conn.RelateArgs,
			RelateExtra: conn.RelateExtra}
	}

	if insert := atom.GetAction("insert"); insert != nil {
		nodeInsert := &Node_Actions_Insert{ActionName: "insert", IsDo: insert.GetIsDo()}
		for _, prepare := range insert.GetPrepares() {
			nodeInsert.PrepareConnects = append(nodeInsert.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range insert.GetNextpages() {
			nodeInsert.NextpageConnects = append(nodeInsert.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.InsertItem = nodeInsert
	}

	if insupd := atom.GetAction("insupd"); insupd != nil {
		nodeInsupd := &Node_Actions_Insupd{ActionName: "insupd", IsDo: insupd.GetIsDo()}
		for _, prepare := range insupd.GetPrepares() {
			nodeInsupd.PrepareConnects = append(nodeInsupd.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range insupd.GetNextpages() {
			nodeInsupd.NextpageConnects = append(nodeInsupd.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.InsupdItem = nodeInsupd
	}

	if iupdate := atom.GetAction("update"); iupdate != nil {
		update := iupdate.(*godbi.Update)
		nodeUpdate := &Node_Actions_Update{
			ActionName: "update",
			IsDo: update.GetIsDo(),
			Empties: update.Empties}
		for _, prepare := range update.GetPrepares() {
			nodeUpdate.PrepareConnects = append(nodeUpdate.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range update.GetNextpages() {
			nodeUpdate.NextpageConnects = append(nodeUpdate.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.UpdateItem = nodeUpdate
	}

	if delett := atom.GetAction("delete"); delett != nil {
		nodeDelete := &Node_Actions_Delete{ActionName: "delete", IsDo: delett.GetIsDo()}
		for _, prepare := range delett.GetPrepares() {
			nodeDelete.PrepareConnects = append(nodeDelete.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range delett.GetNextpages() {
			nodeDelete.NextpageConnects = append(nodeDelete.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.DeleteItem = nodeDelete
	}

	if delecs := atom.GetAction("delecs"); delecs != nil {
		nodeDelecs := &Node_Actions_Delecs{ActionName: "delecs", IsDo: delecs.GetIsDo()}
		for _, prepare := range delecs.GetPrepares() {
			nodeDelecs.PrepareConnects = append(nodeDelecs.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range delecs.GetNextpages() {
			nodeDelecs.NextpageConnects = append(nodeDelecs.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.DelecsItem = nodeDelecs
	}

	if itopics := atom.GetAction("topics"); itopics != nil {
		topics := itopics.(*godbi.Topics)
		nodeTopics := &Node_Actions_Topics{
			IsDo: topics.GetIsDo(),
			ActionName: "topics",
			FIELDS: topics.FIELDS,
			TotalForce: int32(topics.TotalForce),
			MAXPAGENO: topics.MAXPAGENO,
			TOTALNO: topics.TOTALNO,
			ROWCOUNT: topics.ROWCOUNT,
			PAGENO: topics.PAGENO,
			SORTBY: topics.SORTBY,
			SORTREVERSE: topics.SORTREVERSE}
		for _, joint := range topics.Joints {
			nodeTopics.Joints = append(nodeTopics.Joints, &Node_Actions_Joint{
				TableName: joint.TableName,
				Alias: joint.Alias,
				JoinType: joint.JoinType,
				JoinUsing: joint.JoinUsing,
				JoinOn: joint.JoinOn,
				Sortby: joint.Sortby})
		}
		for _, prepare := range topics.GetPrepares() {
			nodeTopics.PrepareConnects = append(nodeTopics.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range topics.GetNextpages() {
			nodeTopics.NextpageConnects = append(nodeTopics.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.TopicsItem = nodeTopics
	}

	if iedit := atom.GetAction("edit"); iedit != nil {
		edit := iedit.(*godbi.Edit)
		nodeEdit := &Node_Actions_Edit{
			IsDo: edit.GetIsDo(),
			ActionName: "edit",
			FIELDS: edit.FIELDS}
		for _, joint := range edit.Joints {
			nodeEdit.Joints = append(nodeEdit.Joints, &Node_Actions_Joint{
				TableName: joint.TableName,
				Alias: joint.Alias,
				JoinType: joint.JoinType,
				JoinUsing: joint.JoinUsing,
				JoinOn: joint.JoinOn,
				Sortby: joint.Sortby})
		}
		for _, prepare := range edit.GetPrepares() {
			nodeEdit.PrepareConnects = append(nodeEdit.PrepareConnects, dbiConnection(prepare))
		}
		for _, nextpage := range edit.GetNextpages() {
			nodeEdit.NextpageConnects = append(nodeEdit.NextpageConnects, dbiConnection(nextpage))
		}
		nodeActions.EditItem = nodeEdit
	}

	return nodeActions
}
