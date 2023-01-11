package gometa

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/genelet/molecule/godbi"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type Restful interface {
	SetPK(string)
	GetPK() string
    Write(context.Context, *sql.DB, proto.Message, ...map[string]interface{}) ([]interface{}, error)
    Read(context.Context, *sql.DB, proto.Message, ...map[string]interface{}) ([]interface{}, error)
    List(context.Context, *sql.DB, proto.Message, ...map[string]interface{}) ([]interface{}, error)
    Update(context.Context, *sql.DB, proto.Message, ...map[string]interface{}) ([]interface{}, error)
    Delete(context.Context, *sql.DB, proto.Message, ...map[string]interface{}) ([]interface{}, error)
}

var _ Restful = (*Rest)(nil)

// Rest handles simple restful actions on Protobuf message,
// using the engine described by Molecule.
//
type Rest struct {
	*godbi.Molecule
	pk   string
}

func NewRest(raw []byte, pk ...string) (*Rest, error) {
	m, err := godbi.NewMoleculeJson(json.RawMessage(raw))
	if err != nil {
		return nil, err
	}

	if pk == nil {
		return &Rest{Molecule:m}, nil
	}
	return &Rest{Molecule:m, pk: pk[0]}, nil
}

func NewGraphRest(graph *Graph) *Rest {
	m, _ := GraphToMolecule(graph)
	return &Rest{Molecule:m, pk: graph.PkName}
}

func (self *Rest) GetPK() string {
	return self.pk
}

func (self *Rest) SetPK(pk string) {
	self.pk = pk
}

func (self *Rest) nameArgsFromPBExtra(check bool, pb proto.Message, extra ...map[string]interface{}) (string, map[string]interface{}, error) {
	name := string(pb.ProtoReflect().Descriptor().Name())
	args, err := ProtobufToMap(pb)
	if err != nil { return "", nil, err }
	if extra != nil {
		for k, v := range extra[0] {
			args[k] = v
		}
	}
	if check {
		if self.pk == "" {
			return "", nil, fmt.Errorf("primary key not defined")
		}
		if _, ok := args[self.pk]; !ok {
			return "", nil, fmt.Errorf("primary key is empty")
		}
	}
	return name, args, nil
}

// Search protobuf messages by placeholder's protobuf definition, with optional constraint extra.
//
func (self *Rest) List(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]interface{}) ([]interface{}, error ) {
	name := string(pb.ProtoReflect().Descriptor().Name())
	if extra != nil {
		return self.Molecule.RunContext(ctx, db, name, "topics", nil, extra[0])
	}
	return self.Molecule.RunContext(ctx, db, name, "topics")
}

// Get proto message from database by the primary key defined in constraint extra.
//
func (self *Rest) Read(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]interface{}) ([]interface{}, error) {
	name, args, err := self.nameArgsFromPBExtra(true, pb, extra...)
    if err != nil {
        return nil, err
    }

	return self.Molecule.RunContext(ctx, db, name, "edit", args)
}

// Insert protobuf message into database, with optional input data in extra.
//
func (self *Rest) Write(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]interface{}) ([]interface{}, error ) {
	name, args, err := self.nameArgsFromPBExtra(false, pb, extra...)
    if err != nil {
        return nil, err
    }
    return self.Molecule.RunContext(ctx, db, name, "insert", args)
}

// Update protobuf message in database by the primary key defined in contraint extra.
//
func (self *Rest) Update(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]interface{}) ([]interface{}, error ) {
	name, args, err := self.nameArgsFromPBExtra(true, pb, extra...)
    if err != nil {
        return nil, err
    }
    return self.Molecule.RunContext(ctx, db, name, "update", args)
}

// Delete protobuf message from database by the primary key defined in constraint extra.
//
func (self *Rest) Delete(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]interface{}) ([]interface{}, error ) {
	name, args, err := self.nameArgsFromPBExtra(true, pb, extra...)
    if err != nil {
        return nil, err
    }
    return self.Molecule.RunContext(ctx, db, name, "delecs", args)
}

func ProtobufToMap(pb proto.Message) (map[string]interface{}, error) {
	m := protojson.MarshalOptions{AllowPartial: true}
    bs, err := m.Marshal(pb)
    if err != nil {
        return nil, err
    }
    hash := make(map[string]interface{})
    err = json.Unmarshal(bs, &hash)
    return hash, err
}

// MapsToProtobufs converts multiple items returned from REST to a slice of protobuf defined in pb.
func MapsToProtobufs(lists []interface{}, pb proto.Message) ([]proto.Message, error) {
	var pbs []proto.Message
	for _, item := range lists {
		newPb := proto.Clone(pb)
		hash, ok := item.(map[string]interface{})
		if !ok { return nil, fmt.Errorf("wrong data type for item: %T", item) }
		err := MapToProtobuf(hash, newPb)
		if err != nil { return nil, err }
		pbs = append(pbs, newPb)
	}
	return pbs, nil
}

// MapToProtobuf converts an item, which is a map, returned from REST to protobuf pb.
func MapToProtobuf(item map[string]interface{}, pb proto.Message) error {
    bs, err := json.Marshal(item)
    if err != nil { return err }
	um := protojson.UnmarshalOptions{DiscardUnknown: true, AllowPartial: true}
    return um.Unmarshal(bs, pb)
}
