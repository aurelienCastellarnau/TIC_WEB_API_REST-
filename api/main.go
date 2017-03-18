package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var claimRequire *sql.DB
var sqlErr error
var claimRequireErr error

func main() {
	db, sqlErr = sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/rest")
	claimRequire, claimRequireErr = sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/claims")
	if claimRequireErr != nil {
		fmt.Printf("\n[claim require] database claims not avalaible. Authentication can't be persisted.")
	}
	if sqlErr != nil {
		fmt.Printf(sqlErr.Error())
		return
	}
	defer db.Close()
	sqlErr = db.Ping()
	if sqlErr != nil {
		fmt.Printf(sqlErr.Error())
		return
	}
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":3000", router))
}
