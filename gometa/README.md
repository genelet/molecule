# molecule/gometa

In _molecule_, a table and its associated actions build up an _atom_. Atoms and relationships between them build up a _molecule_. While traditional REST acts on individual table, _molecule_ acts on whole RDB across all tables.

A molecule could be also represented by protocol buffer, called _Graph_:

[https://github.com/genelet/molecule/blob/master/gometa/meta.proto](https://github.com/genelet/molecule/blob/master/gometa/meta.proto)

This package implements 2 functions to translate _Molecule_ to _Graph_, and
_Graph_ to _Molecule_.

<br /><br />

## Chapter 1. Oneof

_oneof_ is a powerful message type in protobuf, yet it is not associated with
any native GO type. In this package, we associate elements in _oneof_ to a list of table columns. If a table has mutiple _oneof_, they will be represented by string-to-list map, where the key is oneof's name and the value a list of columns.

```go
map[string][]string
```

For a whole database which consists of many tables, all the oneofs are
represented by

```go
map[tableName][string][]string
```

<br /><br />

## Chapter 2. Functions

### 2.1 _Graph_ to _Molecule_

```go
func GraphToMolecule(graph *Graph) (*godbi.Molecule, map[string]map[string][]string) {
```

It translates _Graph_ to _Molecule_ and the associated oneofs in _map[string]map[string][]string_.

<br />

### 2.2 _Molecule_ to _Graph_

```go
func MoleculeToGraph(molecule *godbi.Molecule, rest ...interface{}) *Graph
```

where the first element of _rest_ is the _oneof_ map, and the second _Graph_ package name.


<br />

### 2.3 Errors

Because of the exact matchings between protobuf fields to table columns, there is
no ambiguity in the functions. Every error should be panic and
be fixed in the package.
