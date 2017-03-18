package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Role API definition type (cf consigne)
type Role struct {
	Admin bool `json:",string,admin"`
	User  bool `json:",string,user"`
}

// Phone API definition type (cf consigne)
type Phone struct {
	Number string `json:"number,omitempty"`
}

// User Model
// use to map the database's table user
type User struct {
	UID       int    `json:"uid,string"`
	Lastname  string `json:"lastname,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      Role   `json:"role"`
	Phones    Phone  `json:",omitempty"`
	Password  string `json:"password,omitempty"`
}

// Users slice of several User
type Users []User

// Get fullfill the current Users instance for the response to GET /user
func (users *Users) Get() error {
	role := ""
	rows, sqlErr := db.Query("SELECT * FROM user")
	switch {
	case sqlErr == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case sqlErr != nil:
		log.Printf("%s", sqlErr)
	default:
		fmt.Printf("Users are loaded.")
		for rows.Next() {
			var user User
			rows.Scan(
				&user.UID,
				&user.Firstname,
				&user.Lastname,
				&user.Email,
				&user.Password,
				&role)
			user.ParseRole(role)
			*users = append(*users, user)
		}
	}
	return sqlErr
}

// Assert check user credential and convert a valid credentials to cookie
// Use the Claim interface to manage the retrieving || creation of the token's claim
func (users *Users) Assert(user User) (http.Cookie, error) {
	var cookie http.Cookie
	var err error

	if len(*users) == 0 {
		users.Get()
	}
	for _, u := range *users {
		if user.Email == u.Email && user.Password == u.Password {
			var claim Claim
			user = u
			fmt.Printf("\n%s %s est authentifié avec le role %s.", user.Firstname, user.Lastname, user.Role.toString())
			claim, token, sqlErr := user.GetClaim()
			if sqlErr != nil {
				return cookie, sqlErr
			}
			err := claim.GetValidClaim(user, token)
			if err != nil {
				fmt.Printf("\n[stack trace] <Users.Assert()>")
				return cookie, err
			}
			jwToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			signedToken, _ := jwToken.SignedString([]byte("secret"))
			cookie := http.Cookie{
				Name:    "Auth",
				Value:   signedToken,
				Expires: time.Unix(claim.StandardClaims.ExpiresAt, 0),
			}
			return cookie, nil
		}
	}
	err = errors.New("invalid credentials")
	return cookie, err
}

// GetByID fulfill the current user instance for the response to GET /user/{id}
func (u *User) GetByID() error {
	var role string
	db.QueryRow("SELECT * FROM user WHERE id=?", u.UID).Scan(
		&u.UID,
		&u.Firstname,
		&u.Lastname,
		&u.Email,
		&u.Password,
		&role)
	switch {
	case sqlErr == sql.ErrNoRows:
		fmt.Printf("\nNo user with that ID.")
	case sqlErr != nil:
		fmt.Printf("\n%s", sqlErr.Error())
	default:
		fmt.Printf("\nUser %s %s ID %d\n", u.Firstname, u.Lastname, u.UID)
	}
	u.ParseRole(role)
	return sqlErr
}

// Post send the current user to the database
func (u *User) Post() (int, error) {
	if err := u.Check(); err != nil {
		fmt.Printf("\nuser.check(): invalid user. %s", err.Error())
		return 0, err
	}
	stmt, sqlErr := db.Prepare("INSERT INTO user (lastname, firstname, email, password, role) VALUES (?, ?, ?, ?, ?)")
	if sqlErr != nil {
		return 0, sqlErr
	}
	err := u.GetPasswordHash()
	if err != nil {
		return 0, err
	}
	result, sqlErr := stmt.Exec(u.Lastname, u.Firstname, u.Email, u.Password, u.Role.toString())
	id, sqlErr := result.LastInsertId()
	return int(id), sqlErr
}

// GetClaim retrieve claim from user
func (u *User) GetClaim() (claim Claim, token string, sqlErr error) {
	var database string
	var tmpInt int
	var tmpSource string

	if u.UID == 0 {
		return
	}
	sqlErr = db.QueryRow("SELECT DATABASE()").Scan(&database)
	sqlErr = claimRequire.QueryRow("SELECT * FROM claim WHERE idUser = ? AND source=?", u.UID, database).Scan(
		&tmpInt,
		&claim.ID,
		&tmpSource,
		&token,
	)
	return
}

// Check control user's data
func (u User) Check() error {
	fmt.Println(u)
	switch {
	case u.Firstname == "":
		return errors.New("\nthe firstname isn't valid")
	case u.Lastname == "":
		return errors.New("\nthe lastname isn't valid")
	case u.Email == "":
		return errors.New("\nthe email isn't valid")
	case u.Password == "":
		return errors.New("\nthe password isn't valid")
	case u.Role.toString() == "":
		return errors.New("\nthe role isn't valid")
	}
	return nil
}

// GetPasswordHash transform clear password to sha1 string
func (u *User) GetPasswordHash() error {
	if u.Password == "" {
		err := errors.New("\nno password in current user")
		return err
	}
	pass := sha1.New()
	io.WriteString(pass, u.Password)
	u.Password = hex.EncodeToString(pass.Sum(nil))
	return nil
}

// ParseRole adapt database content to API definition
func (u *User) ParseRole(role string) {
	u.Role.User = true
	if role == "admin" {
		u.Role.Admin = true
	} else {
		u.Role.Admin = false
	}
}

func (r *Role) toString() string {
	if !r.User {
		return ""
	}
	if r.Admin {
		return "admin"
	}
	return "normal"
}

// UnmarshalJSON specific Unmarshal process for Role struct in User
func (r *Role) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		r.Admin = false
		r.User = false
	case "admin":
		r.Admin = true
		r.User = true
	case "normal":
		r.Admin = false
		r.User = true
	}
	return nil
}
