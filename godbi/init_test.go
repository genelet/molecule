package godbi

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func getdb() (*sql.DB, error) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	return sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
}
