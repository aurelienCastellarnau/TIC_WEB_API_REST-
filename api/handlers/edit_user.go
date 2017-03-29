package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"rest/api/helpers"
	"rest/api/model"
	"strconv"

	"github.com/gorilla/mux"
)

// EditUser handler for road PUT /user/{id}
func EditUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.EditUser()")
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
		helpers.SendNotFoundErrorResponse(w, "\n<GET /user/{id}> >> user.GetUserByID(): id < 0 ", nil)
		return
	}
	err = user.GetByID()
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n[EditUser] No user with this UID in database", err)
		fmt.Printf("\nUID: %d", user.UID)
		return
	}
	r.ParseForm()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n[EditUser] Error during Unmarshal body in user object", err)
		return
	}
	id, err := user.Put()
	fmt.Println(user)
	fmt.Printf("\n[LOG DEV] id: %d user id: %d email: %s", id, user.UID, user.Email)
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n[EditUser] Sql update fail: ", err)
		return
	}
	existingUser.Convert(user)
	response, err := json.Marshal(existingUser)
	if err != nil {
		helpers.SendNotFoundErrorResponse(w, "\n[EditUser] Error during serialization into []byte of the existing user.", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Printf("\nUpdate ok. User %d modified.", id)
	fmt.Fprintf(w, "%s", response)
}
