package handlers

import (
	"fmt"
	"net/http"
	"rest/api/helpers"
	"rest/api/model"
	"strconv"

	"rest/api/dao"

	"github.com/gorilla/mux"
)

// DeleteUser handler for road DELETE /user/{id}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.DeleteUser()")
	var user model.User
	var err error

	vars := mux.Vars(r)
	user.UID, err = strconv.Atoi(vars["id"])
	if err != nil || user.UID <= 0 {
		helpers.SendNotFoundErrorResponse(w, "\n<DELETE /user/{id}> >> user.GetUserByID(): fail to parse ID to int: ", err)
		return
	}
	if user.UID <= 0 {
		helpers.SendNotFoundErrorResponse(w, "\n<DELETE /user/{id}> >> user.GetUserByID(): id < 0 ", nil)
		return
	}
	err = user.GetByID()
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n<DELETE /user/{id}> No user with this UID in database", err)
		return
	}
	_, dao.SQLErr = user.Delete()
	if dao.SQLErr != nil {
		helpers.SendNotFoundErrorResponse(w, "\n<DELETE /user/{id}> >> Sql DELETE failed: ", dao.SQLErr)
	}
	w.WriteHeader(http.StatusNoContent)
}
