package main

import (
	"fmt"
	"net/http"

	"errors"

	"strings"

	"encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"
)

// SecureAPIWithToken process the cookie to check protected roads.
func SecureAPIWithToken(process http.Handler, needToBeTokenProtected bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n[Handlers stack trace] call of Handler.SecureAPIWithToken()")
		var user User
		var claim Claim

		if needToBeTokenProtected {
			fmt.Printf("\n\n[SecureAPI stack trace] needToBeProtected. \n")
			tokenCookie, err := r.Cookie("Auth")
			if err != nil || tokenCookie == nil {
				fmt.Printf("\nNo token present.")
				fmt.Fprintf(w, "Must be connected")
				http.Redirect(w, r, "/", http.StatusUnauthorized)
				return
			}
			fmt.Printf("\n\n[SecureAPI stack trace] pass the cookie check. \n")
			token, err := jwt.ParseWithClaims(tokenCookie.String(), &claim, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					err := errors.New("token got unexpected crypting algorythm")
					return nil, err
				}
				return []byte("secret"), nil
			})
			if err != nil {
				fmt.Fprintf(w, "Error response.")
				fmt.Printf("\n<SecureAPI> fail to parse token: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Printf("\n\n[SecureAPI stack trace] token parsed \n")
			claim, ok := token.Claims.(*Claim)
			if ok && token.Valid {
				user.UID = claim.ID
				sqlErr = user.GetByID()
				if sqlErr != nil {
					fmt.Printf("\n<SecureAPI> fail to load user: %s", sqlErr.Error())
					fmt.Fprintf(w, "Internal server error: database request failed.")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if user.Role.toString() != "admin" {
					fmt.Printf("\n<SecureAPI>: your role %s doesn't", user.Role)
					fmt.Fprintf(w, "Must be admin")
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
		}
		process.ServeHTTP(w, r)
	})
}

// SecureAPIWithBasic Protect route with basic http authentication
func SecureAPIWithBasic(process http.Handler, needToBeBasicProtected bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		var users Users

		fmt.Printf("\n[Handlers stack trace] call of Handler.SecureAPIWithBasic()")

		authError := func() {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Must be admin\"")
			http.Error(w, "authorization failed", http.StatusUnauthorized)
		}

		if needToBeBasicProtected {
			fmt.Printf("\n[basic security] This road need BasicHttpAuthentication.")

			auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(auth) != 2 || auth[0] != "Basic" {
				fmt.Printf("\n[basic security] Header Authorization do not permit to identify BasicHttpAuthentication")
				fmt.Println(auth)
				authError()
				return
			}
			decodedCredentials, err := base64.StdEncoding.DecodeString(auth[1])
			if err != nil {
				fmt.Printf("\n[basic security] base64 encoded 'Authorization' content not properly decoded.")
				fmt.Printf("\n %s", auth[1])
				authError()
				return
			}
			credentials := strings.SplitN(string(decodedCredentials), ":", 2)
			if len(credentials) == 2 {
				user.Email = credentials[0]
				user.Password = credentials[1]
				err = user.GetPasswordHash()
				_, err = users.Assert(user)
			}
			if err != nil {
				fmt.Printf("\n[basic security] user not asserted.")
				fmt.Printf("\n email: %s, Password: %s", credentials[0], credentials[1])
				authError()
				return
			}
		}
		process.ServeHTTP(w, r)
	})
}
