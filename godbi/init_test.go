package godbi

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func getdb() (*sql.DB, error) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := "gotest"
	if dbUser == "" { return nil, fmt.Errorf("missing DBUSER") }
	if dbPass == "" { return nil, fmt.Errorf("missing DBPASS") }
	if dbName == "" { return nil, fmt.Errorf("missing DBNAME") }
	dbHOST := os.Getenv("DBHOST")
 	if dbHOST != "" {
		return sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHOST+")/"+dbName)
	}
	return sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
}
