package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/gorilla/mux"
)

// TokenResponse structure to serialize cookie as json
type TokenResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// ErrorResponse API definition type (cf consigne)
type ErrorResponse struct {
	Code          int    `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	Missingfields string `json:"missingFields,omitempty"`
}

// Index handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	http.FileServer(http.Dir("./"))
}

// SetToken Handler to set the athentication token
func SetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.SetToken()")
	var user User
	var users Users
	var auth TokenResponse

	r.ParseForm()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("\n<POST /token> >> json.Unmarshal(body, &user): %s", err.Error())
		http.Error(w, "ErrorResponse", http.StatusBadRequest)
		return
	}
	if err = user.GetPasswordHash(); err != nil {
		fmt.Printf("\n<POST /token> >> user.GetPasswordHash(): %s", err.Error())
		http.Error(w, "ErrorResponse", http.StatusBadRequest)
		return
	}
	if sqlErr = users.Get(); sqlErr != nil {
		fmt.Printf("\n<POST /token> >> users.Get(): %s", sqlErr.Error())
		http.Error(w, "ErrorResponse", http.StatusBadRequest)
		return
	}
	cookie, err := users.Assert(user)
	if err != nil {
		fmt.Printf("\n<POST /token> >> users.Assert(): %s", err.Error())
		http.Error(w, "Must be connected - "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Printf("\n[Handlers stack trace] Token set redirection with status 200")
	auth.Name = "Auth"
	auth.Token = cookie.Value
	response, err := json.Marshal(auth)
	http.SetCookie(w, &cookie)
	w.Write(response)
	http.Redirect(w, r, "/", 200)
}

//GetUsers handler
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.GetUsers()")
	var users Users

	sqlErr = users.Get()
	if sqlErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "ErrorResponse")
		fmt.Printf("\n<GET /user> >> user.GetUsers(): fail to get users %s", sqlErr.Error())
	}
	response, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<GET /user> >> user.GetUsers(): fail to Marshall users %s", err.Error())
	}
	fmt.Fprintf(w, "%s", response)
}

// GetUserByID handler
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.GetUserById()")
	var user User
	var err error

	vars := mux.Vars(r)
	user.UID, err = strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<GET /user/{id}> >> user.GetUserByID(): fail to parse ID to int %s", err.Error())
	}
	sqlErr = user.GetByID()
	if sqlErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<GET /user/{id}> >> user.GetUserByID(): fail to get user %s", sqlErr.Error())
		return
	}
	response, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<GET /user/{id}> >> user.GetUserByID(): fail to Marshall users %s", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", response)
}

// AddUser handler
func AddUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[Handlers stack trace] call of Handler.AddUser()")
	var user User
	var users Users

	r.ParseForm()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println(body)
	json.Unmarshal(body, &user)

	fmt.Printf("CHECK USER: firstname: %s, lastname: %s, email: %s, password: %s, role: %s", user.Firstname, user.Lastname, user.Email, user.Password, user.Role.Admin)
	id, sqlErr := user.Post()
	if sqlErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<POST /user> >> user.adduser(): fail to persist %s", sqlErr.Error())
		return
	}
	sqlErr = users.Get()
	if sqlErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<POST /user> >> user.adduser(): fail to load users %s", sqlErr.Error())
		return
	}
	response, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ErrorResponse")
		fmt.Printf("\n<POST /user> >> user.adduser(): fail to Marshall users %s", err.Error())
		return
	}
	fmt.Printf("\nInsert ok. New user created with id: %d", id)
	fmt.Fprintf(w, "%s", response)
	w.WriteHeader(http.StatusCreated)
}
