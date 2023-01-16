# molecule/gometa

In _molecule_, a table and its associated actions build up an _atom_. Atoms and relationships between them build up a _molecule_. While traditional REST acts on individual table, _molecule_ acts on whole RDB across all tables.

A molecule could be also represented by protocol buffer, called _Graph_:

[https://github.com/genelet/molecule/blob/master/gometa/meta.proto](https://github.com/genelet/molecule/blob/master/gometa/meta.proto)

This package implements 2 functions to translate _Molecule_ to _Graph_, and
_Graph_ to _Molecule_.

<br /><br />

## Chapter 1. Oneof

_oneof_ is a powerful message type in protobuf, yet it is not associated with
any native GO type. In this package, we assign _oneof_ to a list of table
columns. If table has mutiple _oneof_, they will be represented by a map,
in which key is the name of oneof and value its list of columns.

```go
map[string][]string
```

For whole database which contains many tables, all the oneofs are
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
func MoleculeToGraph(molecule *godbi.Molecule, args ...interface{}) *Graph
```

where _args_ are
  - oneofs in the database: map[atomName][oneofName][]string
  - package name
  - goPackage name
  - primary table name
  - primary table's pk
  - child to parent table mapping
  - table name to pk mapping

<br />

### 2.3 Errors

Matchings between protobuf fields to table columns have to be unique, hence,
there is no ambiguity in the functions. Every error should be panic and
be fixed in the package.
