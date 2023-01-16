# molecule/rdb

In _molecule_, a table and its associated actions build up an _atom_. Atoms and relationships between them build up a _molecule_. While traditional REST acts on individual table, _molecule_ acts on whole RDB across all tables.

<br /><br />

## Molecule Creation 

This package generates automatically a simple molecule from relational
database, using meta data such as primary keys, foreign keys, and 
auto increment fields etc.

Currently only the following three database types are supported:
_godbi.Postgres_, _godbi.MySQL_ and _godbi.SQLite_.

```go
func NewMolecule(db *sql.DB, driver godbi.DBType, dbName string) (*godbi.Molecule, error)
```

where _db_ is the standard database handler, _driver_ the [DBType](https://github.com/genelet/molecule#chapter-5-molecule-usage) defined in _Molecule_ and _dbName_ the name of database.

After _molecule_ is created, we can run RESTful actions on multiple tables at
once.
