package main

import (
	"log"
	"net/http"

	"rest/api/routing"

	"rest/api/dao"
)

func main() {
	dao.GetDb()
	defer dao.Db.Close()
	defer dao.ClaimRequire.Close()
	router := routing.NewRouter()
	log.Fatal(http.ListenAndServe(":3000", router))
}
