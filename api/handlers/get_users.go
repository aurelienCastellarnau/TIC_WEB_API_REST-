package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest/api/dao"
	"rest/api/helpers"
	"rest/api/model"
)

//GetUsers handler for road GET /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.GetUsers()")
	var users model.Users
	var existingUsers model.ExistingUsers

	dao.SQLErr = users.Get()
	if dao.SQLErr != nil {
		helpers.SendErrorResponse(w, "\n<GET /user> >> user.GetUsers(): fail to get users: ", dao.SQLErr)
		return
	}
	existingUsers.Convert(users)
	response, err := json.Marshal(existingUsers)
	if err != nil {
		helpers.SendErrorResponse(w, "\n<GET /user> >> user.GetUsers(): fail to Marshall users: ", err)
	}
	fmt.Fprintf(w, "%s", response)
}
