package gometa

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/genelet/molecule/godbi"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Restful interface {
	Write(context.Context, *sql.DB, proto.Message, ...map[string]any) ([]any, error)
	Read(context.Context, *sql.DB, proto.Message, ...map[string]any) ([]any, error)
	List(context.Context, *sql.DB, proto.Message, ...map[string]any) ([]any, error)
	Update(context.Context, *sql.DB, proto.Message, ...map[string]any) ([]any, error)
	Delete(context.Context, *sql.DB, proto.Message, ...map[string]any) ([]any, error)
}

var _ Restful = (*Rest)(nil)

// Rest handles simple restful actions on Protobuf message,
// using the engine described by Molecule.
type Rest struct {
	*Graph
	mole *godbi.Molecule
}

func NewRest(graph *Graph) *Rest {
	mole, _ := GraphToMolecule(graph)
	return &Rest{graph, mole}
}

func NewRestByte(bs []byte) (*Rest, error) {
	graph := new(Graph)
	um := protojson.UnmarshalOptions{DiscardUnknown: true, AllowPartial: true}
	err := um.Unmarshal(bs, graph)
	if err != nil {
		return nil, err
	}

	return NewRest(graph), nil
}

func (r *Rest) nameArgsFromPBExtra(check bool, pb proto.Message, extra ...map[string]any) (string, map[string]any, error) {
	name := string(pb.ProtoReflect().Descriptor().Name())
	args, err := ProtobufToMap(pb)
	if err != nil {
		return "", nil, err
	}
	if extra != nil {
		for k, v := range extra[0] {
			args[k] = v
		}
	}
	if check {
		if r.Graph.PkName == "" {
			return "", nil, fmt.Errorf("primary key not defined")
		}
		if _, ok := args[r.Graph.PkName]; !ok {
			return "", nil, fmt.Errorf("the primary key, %s is empty", r.Graph.PkName)
		}
	}
	return name, args, nil
}

// Search protobuf messages by placeholder's protobuf definition, with optional constraint extra.
func (r *Rest) List(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]any) ([]any, error) {
	name := string(pb.ProtoReflect().Descriptor().Name())
	if extra != nil {
		return r.mole.RunContext(ctx, db, name, "topics", &godbi.RunOption{Extra: extra[0]})
	}
	return r.mole.RunContext(ctx, db, name, "topics", nil)
}

// Get proto message from database by the primary key defined in constraint extra.
func (r *Rest) Read(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]any) ([]any, error) {
	name, args, err := r.nameArgsFromPBExtra(true, pb, extra...)
	if err != nil {
		return nil, err
	}

	return r.mole.RunContext(ctx, db, name, "edit", &godbi.RunOption{Args: args})
}

// Insert protobuf message into database, with optional input data in extra.
func (r *Rest) Write(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]any) ([]any, error) {
	name, args, err := r.nameArgsFromPBExtra(false, pb, extra...)
	if err != nil {
		return nil, err
	}
	return r.mole.RunContext(ctx, db, name, "insert", &godbi.RunOption{Args: args})
}

// Update protobuf message in database by the primary key defined in contraint extra.
func (r *Rest) Update(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]any) ([]any, error) {
	name, args, err := r.nameArgsFromPBExtra(true, pb, extra...)
	if err != nil {
		return nil, err
	}
	return r.mole.RunContext(ctx, db, name, "update", &godbi.RunOption{Args: args})
}

// Delete protobuf message from database by the primary key defined in constraint extra.
func (r *Rest) Delete(ctx context.Context, db *sql.DB, pb proto.Message, extra ...map[string]any) ([]any, error) {
	name, args, err := r.nameArgsFromPBExtra(true, pb, extra...)
	if err != nil {
		return nil, err
	}
	return r.mole.RunContext(ctx, db, name, "delecs", &godbi.RunOption{Args: args})
}

func ProtobufToMap(pb proto.Message) (map[string]any, error) {
	m := protojson.MarshalOptions{AllowPartial: true}
	bs, err := m.Marshal(pb)
	if err != nil {
		return nil, err
	}
	hash := make(map[string]any)
	err = json.Unmarshal(bs, &hash)
	return hash, err
}

// MapsToProtobufs converts multiple items returned from REST to a slice of protobuf defined in pb.
func MapsToProtobufs(lists []any, pb proto.Message) ([]proto.Message, error) {
	var pbs []proto.Message
	for _, item := range lists {
		newPb := proto.Clone(pb)
		hash, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("wrong data type for item: %T", item)
		}
		err := MapToProtobuf(hash, newPb)
		if err != nil {
			return nil, err
		}
		pbs = append(pbs, newPb)
	}
	return pbs, nil
}

// MapToProtobuf converts an item, which is a map, returned from REST to protobuf pb.
func MapToProtobuf(item map[string]any, pb proto.Message) error {
	bs, err := json.Marshal(item)
	if err != nil {
		return err
	}
	um := protojson.UnmarshalOptions{DiscardUnknown: true, AllowPartial: true}
	return um.Unmarshal(bs, pb)
}
