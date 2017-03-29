package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"rest/api/dao"
	"strings"

	"rest/api/helpers"
	"strconv"
)

// Role API definition type (cf consigne)
type Role struct {
	Admin bool `json:",omitempty,string,admin"`
	User  bool `json:",omitempty,string,user"`
}

// Phone API definition type (cf consigne)
type Phone struct {
	Number string `json:"number,omitempty"`
}

// ExistingUser API definition for user (cf: consigne)
type ExistingUser struct {
	UID       int    `json:"uid,string,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      Role   `json:"role,omitempty"`
	Phones    Phones `json:"phones,omitempty"`
}

// User Model
// use to map the database's table user
type User struct {
	UID       int    `json:"uid,string,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      Role   `json:"role,omitempty"`
	Phones    Phones `json:"phones,omitempty"`
	Password  string `json:"password,omitempty"`
}

// Users slice of several User
type Users []User

// ExistingUsers slice of several ExistingUser
type ExistingUsers []ExistingUser

// Phones slice of several Phone
type Phones []Phone

// CredentialUsers slice of several CredentialUser
// type CredentialUsers []CredentialUser

// Convert User in ExistingUser
func (ExistingUser *ExistingUser) Convert(user User) {
	if user.UID == 0 {
		return
	}
	ExistingUser.UID = user.UID
	ExistingUser.Firstname = user.Firstname
	ExistingUser.Lastname = user.Lastname
	ExistingUser.Email = user.Email
	ExistingUser.Phones = user.Phones
	ExistingUser.Role = user.Role
}

// Convert Users in ExistingUsers
func (ExistingUsers *ExistingUsers) Convert(users Users) {
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		var ExistingUser ExistingUser

		ExistingUser.Convert(user)
		*ExistingUsers = append(*ExistingUsers, ExistingUser)
	}
}

// Get fullfill the current Users instance for the response to GET /users
func (users *Users) Get() error {
	var rows *sql.Rows
	var role string

	rows, dao.SQLErr = dao.Db.Query("SELECT * FROM user")
	switch {
	case dao.SQLErr == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case dao.SQLErr != nil:
		log.Printf("%s", dao.SQLErr)
	default:
		for rows.Next() {
			var user User
			rows.Scan(
				&user.UID,
				&user.Lastname,
				&user.Firstname,
				&user.Email,
				&user.Password,
				&role)
			user.ParseRole(role)
			*users = append(*users, user)
		}
		fmt.Printf("\nUsers are loaded.")
	}
	return dao.SQLErr
}

// Assert check user credential
// Possibility to use the Claim interface to manage the retrieving || creation of the token's claim
func (users *Users) Assert(user *User) error {
	var err error

	if len(*users) == 0 {
		users.Get()
	}
	for _, u := range *users {
		if user.Email == u.Email && user.Password == u.Password {
			*user = u
			fmt.Printf("\n%s %s est authentifié avec le role %s.", user.Firstname, user.Lastname, user.Role.ToString())
			return nil
		}
	}
	err = errors.New("invalid credentials")
	return err
}

// Search perform Select Like on name and email
func (users *Users) Search(input string, count int) error {
	var rows *sql.Rows
	var role string
	if count > 0 {
		rows, dao.SQLErr = dao.Db.Query("SELECT * FROM user WHERE firstname LIKE '%" + input + "%' OR lastname LIKE '%" + input + "%' OR email LIKE '%" + input + "%' LIMIT " + strconv.Itoa(count))
	} else {
		rows, dao.SQLErr = dao.Db.Query("SELECT * FROM user WHERE firstname LIKE '%" + input + "%' OR lastname LIKE '%" + input + "%' OR email LIKE '%" + input + "%'")
	}
	switch {
	case dao.SQLErr == sql.ErrNoRows:
		fmt.Printf("No user with that ID.")
	case dao.SQLErr != nil:
		fmt.Printf("%s", dao.SQLErr)
	default:
		for rows.Next() {
			var user User
			rows.Scan(
				&user.UID,
				&user.Lastname,
				&user.Firstname,
				&user.Email,
				&user.Password,
				&role)
			user.ParseRole(role)
			*users = append(*users, user)
		}
		fmt.Printf("\nUsers are loaded.")
	}
	return dao.SQLErr
}

// GetByID fulfill the current user instance for the response to GET /user/{id}
func (u *User) GetByID() error {
	var role string
	dao.Db.QueryRow("SELECT * FROM user WHERE id=?", u.UID).Scan(
		&u.UID,
		&u.Firstname,
		&u.Lastname,
		&u.Email,
		&u.Password,
		&role)
	if u.Email == "" {
		return errors.New("no user at this id")
	}
	u.ParseRole(role)
	return dao.SQLErr
}

// Post send the current user to the database
func (u *User) Post() (int, error) {
	var stmt *sql.Stmt
	var result sql.Result
	var err error
	var id int64

	if err := u.Check(); err != nil {
		fmt.Printf("\nuser.check(): invalid user. %s", err.Error())
		return 0, err
	}
	stmt, dao.SQLErr = dao.Db.Prepare("INSERT INTO user (lastname, firstname, email, password, role) VALUES (?, ?, ?, ?, ?)")
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	defer stmt.Close()
	u.Password, err = helpers.ToSha1(u.Password)
	if err != nil {
		return 0, err
	}
	result, dao.SQLErr = stmt.Exec(u.Lastname, u.Firstname, u.Email, u.Password, u.Role.ToString())
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	id, dao.SQLErr = result.LastInsertId()
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	return int(id), nil
}

// Put modify the designated user in database
func (u *User) Put() (int, error) {
	var set string
	var stmt *sql.Stmt
	var result sql.Result
	var id int64

	properties, err := u.GetDiff()
	if err != nil {
		return 0, err
	}
	query := "UPDATE user SET "
	for key := range properties {
		set = set + key + "='" + properties[key] + "', "
	}
	set = set[0 : len(set)-2]
	query = query + set + " WHERE id = ?"
	fmt.Printf("\nQuery : %s", query)
	stmt, dao.SQLErr = dao.Db.Prepare(query)
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	defer stmt.Close()
	u.Password, err = helpers.ToSha1(u.Password)
	if err != nil {
		return 0, err
	}
	result, dao.SQLErr = stmt.Exec(u.UID)
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	u.GetByID()
	return int(id), nil
}

// Delete the User in database, accordingly to UID
func (u *User) Delete() (int, error) {
	var stmt *sql.Stmt
	var result sql.Result
	var id int64

	stmt, dao.SQLErr = dao.Db.Prepare("DELETE FROM user WHERE id=?")
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	result, dao.SQLErr = stmt.Exec(u.UID)
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	id, dao.SQLErr = result.LastInsertId()
	if dao.SQLErr != nil {
		return 0, dao.SQLErr
	}
	fmt.Printf("\n[deleteUSer]Check on LastInsertId: %d", id)
	return 1, nil
}

// GetClaim retrieve claim from user
func (u *User) GetClaim() (claim Claim, token string, sqlErr error) {
	var database string
	var tmpInt int
	var tmpSource string

	if u.UID == 0 {
		return
	}
	dao.SQLErr = dao.Db.QueryRow("SELECT DATABASE()").Scan(&database)
	dao.SQLErr = dao.ClaimRequire.QueryRow("SELECT * FROM claim WHERE idUser = ? AND source=?", u.UID, database).Scan(
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
	case u.Role.ToString() == "":
		return errors.New("\nthe role isn't valid")
	}
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

// GetDiff return slice of string
// representing the properties which are different between the two users
func (u *User) GetDiff() (map[string]string, error) {
	var compare User
	properties := make(map[string]string)

	compare.UID = u.UID
	err := compare.GetByID()
	if err != nil {
		return nil, err
	}
	if u.Firstname != "" && u.Firstname != compare.Firstname {
		properties["firstname"] = u.Firstname
	} else {
		properties["firstname"] = compare.Firstname
		u.Firstname = compare.Firstname
	}
	if u.Lastname != "" && u.Lastname != compare.Lastname {
		properties["lastname"] = u.Lastname
	} else {
		properties["lastname"] = compare.Lastname
		u.Lastname = compare.Lastname
	}
	if u.Email != "" && u.Email != compare.Email {
		properties["email"] = u.Email
	} else {
		properties["email"] = compare.Email
		u.Email = compare.Email
	}
	if u.Password != "" && u.Password != compare.Password {
		properties["password"] = u.Password
	} else {
		properties["password"] = compare.Password
		u.Password = compare.Password
	}
	if u.Role.ToString() != "" && u.Role.ToString() != compare.Role.ToString() {
		properties["role"] = u.Role.ToString()
	} else {
		properties["role"] = compare.Role.ToString()
		u.Role = compare.Role
	}
	return properties, nil
}

// ToString method to stringify boolean struct Role
func (r *Role) ToString() string {
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

// MarshalJSON specific marshal process for Role struct in User
/*
func (r Role) MarshalJSON() ([]byte, error) {
	if !r.User {
		return json.Marshal("")
	}
	if r.Admin {
		return json.Marshal("admin")
	}
	return json.Marshal("normal")
}
*/
