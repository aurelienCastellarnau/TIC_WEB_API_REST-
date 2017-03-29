package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rest/api/dao"
	"rest/api/helpers"
	"rest/api/model"
)

// AddUser handler for road POST /user
// Empty since we just need to send str...
func AddUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.AddUser()")
	var user model.User
	//var users model.Users
	var id int

	r.ParseForm()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal(body, &user)
	id, dao.SQLErr = user.Post()
	if dao.SQLErr != nil {
		helpers.SendErrorResponse(w, "\n<POST /user> >> user.adduser(): fail to persist", dao.SQLErr)
		return
	}
	/*
		dao.SQLErr = users.Get()
		if dao.SQLErr != nil {
			helpers.SendErrorResponse(w, "\n<POST /user> >> user.adduser(): fail to load users", dao.SQLErr)
			return
		}
		response, err := json.Marshal(users)
		if err != nil {
			helpers.SendErrorResponse(w, "\n<POST /user> >> user.adduser(): fail to Marshall users", err)
			return
		}
	*/
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("\nInsert ok. New user created with id: %d", id)
	fmt.Fprintf(w, "%s", "Created User")
}
