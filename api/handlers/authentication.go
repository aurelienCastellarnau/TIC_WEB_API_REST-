package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest/api/helpers"
	"rest/api/model"
)

// TokenResponse structure to serialize cookie as json
type TokenResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// RefreshAuth just send 401 to make the credentials form pop on client side
// and reset credentials.
func RefreshAuth(w http.ResponseWriter, r *http.Request) {
	helpers.AuthError(w)
}

// Auth 'ask Http Basic..' road send back the user credentials on client.(user is loaded by parseAuthorization)
func Auth(w http.ResponseWriter, r *http.Request) {
	var user model.User
	users := &model.Users{}
	var credentials []string
	var err error

	fmt.Printf("\n[Handlers stack trace] call of Handler.Auth()")

	credentials, err = helpers.ParseAuthorization(r)
	if err != nil {
		helpers.SendMustBeConnectedResponse(w, "\n[Auth] error on parsing authorization: ", err)
		return
	}
	user.Email = credentials[0]
	user.Password = credentials[1]
	err = users.Assert(&user)
	if err != nil {
		helpers.SendMustBeConnectedResponse(w, "\n[Auth] error on asserting user by request content or Authorization header: ", err)
		return
	}
	response, err := json.Marshal(user)
	if err != nil {
		helpers.SendMustBeConnectedResponse(w, "\n[Auth] error on marshalling credentials to response: ", err)
		return
	}
	fmt.Fprintf(w, "%s", response)
}
