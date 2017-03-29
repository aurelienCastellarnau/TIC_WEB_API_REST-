package security

import (
	"fmt"
	"net/http"

	"errors"
	"rest/api/helpers"
	"rest/api/model"

	"rest/api/dao"

	jwt "github.com/dgrijalva/jwt-go"
)

// SecureAPIWithToken process the cookie to check protected roads.
// Onbuild, useless for school expectations
func SecureAPIWithToken(process http.Handler, needToBeTokenProtected bool, reserved bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n[Handlers stack trace] call of Handler.SecureAPIWithToken()")
		var user model.User
		var claim model.Claim

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
			claim, ok := token.Claims.(*model.Claim)
			if ok && token.Valid {
				user.UID = claim.ID
				dao.SQLErr = user.GetByID()
				if dao.SQLErr != nil {
					fmt.Printf("\n<SecureAPI> fail to load user: %s", dao.SQLErr.Error())
					fmt.Fprintf(w, "Internal server error: database request failed.")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if user.Role.ToString() != "admin" {
					fmt.Printf("\n<SecureAPI>: your role %s doesn't", user.Role.ToString())
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
func SecureAPIWithBasic(process http.Handler, needToBeBasicProtected bool, reserved bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		users := &model.Users{}

		fmt.Printf("\n[Handlers stack trace] call of Handler.SecureAPIWithBasic()")
		if needToBeBasicProtected {
			fmt.Printf("\n[basic security] This road need BasicHttpAuthentication.")
			credentials, err := helpers.ParseAuthorization(r)
			if err != nil {
				helpers.SendUnauthorisedResponse(&w, "\n[basic security] error on parsing authorization: ", err)
				return
			}
			user.Email = credentials[0]
			user.Password = credentials[1]
			err = users.Assert(&user)
			if err != nil {
				fmt.Printf("\n[basic security] user not asserted: %s.", err.Error())
				fmt.Printf("\n email: %s, Password: %s", user.Email, user.Password)
				helpers.AuthError(w)
				return
			}
			if user.UID == 0 {
				fmt.Printf("\n[basic security] no user asserted after reading authorization headers.")
				helpers.AuthError(w)
			}
			fmt.Printf("\n[basic security] user firstname: %s, role: %s", user.Firstname, user.Role.ToString())
			if reserved && !user.Role.Admin {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"Must be admin\"")
				r.Header.Set("Authorization", "")
				http.Error(w, "Must be admin", http.StatusForbidden)
				return
			}
		}
		process.ServeHTTP(w, r)
	})
}
