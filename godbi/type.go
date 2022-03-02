package godbi

import (
)

type DBType int64

const (
	SQLDefault DBType = iota
	SQLRaw
	SQLite
	MySQL
	Postgres
	TSMillisecond
	TSMicrosecond
)

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
