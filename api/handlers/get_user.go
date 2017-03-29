package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest/api/dao"
	"rest/api/model"
	"strconv"

	"rest/api/helpers"

	"github.com/gorilla/mux"
)

// GetUserByID handler for road GET /user/{id}
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.GetUserById()")
	var user model.User
	var existingUser model.ExistingUser
	var err error

	vars := mux.Vars(r)
	user.UID, err = strconv.Atoi(vars["id"])
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n<GET /user/{id}> >> user.GetUserByID(): fail to parse ID to int: ", err)
		return
	}
	if user.UID <= 0 {
		helpers.SendNotFoundErrorResponse(w, "\n<GET /user/{id}> >> user.GetUserByID(): id <= 0: ", nil)
		return
	}
	dao.SQLErr = user.GetByID()
	if dao.SQLErr != nil {
		helpers.SendNotFoundErrorResponse(w, "\n<GET /user/{id}> >> user.GetUserByID(): fail to get user: ", dao.SQLErr)
		return
	}
	existingUser.Convert(user)
	response, err := json.Marshal(existingUser)
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n<GET /user/{id}> >> user.GetUserByID(): fail to Marshall users: ", err)
		return
	}
	fmt.Fprintf(w, "%s", response)
}
