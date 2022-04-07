# molecule/rdb

In _molecule_, a table and its associated actions build up an _atom_. Atoms and relationships between them build up a _molecule_. While traditional REST acts on individual table, _molecule_ acts on whole RDB across all tables.

This package generates a simple molecule from a relational database, using its
meta such as primary keys, foreign keys, and auto increment fields etc.

<br /><br />

## Chapter 1. Molecule Creation 

```go
func NewMolecule(db *sql.DB, driver godbi.DBType, dbName string) (*godbi.Molecule, error)
```

where _db_ is the standard database handler, _driver_ the [DBType](https://github.com/genelet/molecule#chapter-5-molecule-usage) defined in _Molecule_ and _dbName_ the name of database.

 Currently only the following three database types are supported:
_godbi.Postgres_, _godbi.MySQL_ and _godbi.SQLite_.
