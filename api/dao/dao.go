package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Db instance of database/sql interface
var Db *sql.DB

// ClaimRequire instance for claim database only
var ClaimRequire *sql.DB

// SQLErr error globale to manage sql return
var SQLErr error

// GetDb instanciate databasem/sql interface
func GetDb() {
	var claimRequireErr error

	Db, SQLErr = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/rest")
	ClaimRequire, claimRequireErr = sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/claims")
	if claimRequireErr != nil || ClaimRequire == nil {
		fmt.Printf("\n[claim require] database claims not avalaible. Authentication by token can't be persisted.")
	}
	if SQLErr != nil {
		fmt.Printf("\n[database connection] database not avalaible: %s", SQLErr.Error())
		fmt.Printf("\nThis api work with Mysql, port 3306, database rest, user admin/admin.")
		return
	}
	SQLErr = Db.Ping()
	if SQLErr != nil {
		fmt.Printf("\n[database connection] PING database not avalaible: %s", SQLErr.Error())
		fmt.Printf("\nThis api work with Mysql, port 3306, database rest, user admin/admin.")
		return
	}
	claimRequireErr = ClaimRequire.Ping()
	if claimRequireErr != nil {
		fmt.Printf("\n[claim require] database claims not avalaible. Authentication by token can't be persisted.")
		ClaimRequire.Close()
	}
}
