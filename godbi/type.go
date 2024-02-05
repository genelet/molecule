package godbi

type DBType int32

const (
	SQLDefault DBType = iota
	SQLRaw
	SQLite
	MySQL
	Postgres
	TSNano
)

// LowerName returns the lower case name of the DBType
func (self DBType) LowerName() string {
	switch self {
	case SQLite:
		return "sqlite3"
	case MySQL:
		return "mysql"
	case Postgres:
		return "postgres"
	default:
	}
	return ""
}

// DBTypeByName returns the DBType by name
func DBTypeByName(str string) DBType {
	switch str {
	case "sqlite3", "sqlite", "SQLite3", "SQLite":
		return SQLite
	case "MySQL", "mysql":
		return MySQL
	case "Postgres", "postgres", "PostgreSQL", "postgresql":
		return Postgres
	default:
	}
	return SQLDefault
}
