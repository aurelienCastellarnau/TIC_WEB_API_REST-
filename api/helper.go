package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func authError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Must be connected\"")
	http.Error(w, "Must be connected", http.StatusUnauthorized)
}

func parseAuthorization(user *User, r *http.Request) ([]string, error) {
	var users Users
	var err error

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		fmt.Printf("\n[basic security] Header Authorization do not permit to identify BasicHttpAuthentication.\n")
		fmt.Println(auth)
		return nil, nil
	}
	decodedCredentials, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		fmt.Printf("\n[basic security] base64 encoded 'Authorization' content not properly decoded.")
		fmt.Printf("\n %s", auth[1])
		return nil, err
	}
	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) == 2 {
		user.Email = credentials[0]
		user.Password = credentials[1]
		err = user.GetPasswordHash()
		_, err = users.Assert(user)
	}
	return credentials, err
}
