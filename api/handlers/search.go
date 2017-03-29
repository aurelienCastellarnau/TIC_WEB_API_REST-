package handlers

import (
	"fmt"
	"net/http"
	"rest/api/model"
	"strconv"

	"rest/api/dao"

	"encoding/json"

	"errors"

	"github.com/gorilla/mux"
)

// InlineResponse200 to fit Api's definition of TIC/WEB Etna 2019
type InlineResponse200 struct {
	Data model.ExistingUsers `json:"data"`
}

// Search retrieve users from database by name or email
func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.Search()")
	var input string
	var count int
	var users model.Users
	var existingUsers model.ExistingUsers
	var inlineResponse200 InlineResponse200
	var err error

	vars := mux.Vars(r)
	input = vars["q"]
	if input == "" {
		err = errors.New("input empty")
		return
	}
	count, err = strconv.Atoi(vars["count"])
	if err != nil {
		err = errors.New("error converting count")
		return
	}
	dao.SQLErr = users.Search(input, count)
	if dao.SQLErr != nil {
		err = errors.New("sql fail")
	}
	existingUsers.Convert(users)
	inlineResponse200.Data = existingUsers
	response, err := json.Marshal(inlineResponse200)
	if err == nil && len(existingUsers) > 0 {
		w.WriteHeader(200)
		w.Write(response)
		return
	}
	w.WriteHeader(200)
}
