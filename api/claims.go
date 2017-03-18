package main

import (
	"encoding/json"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claim for jwt
type Claim struct {
	ID int `json:"ID"`
	jwt.StandardClaims
}

// Claims slice to stock all claims
type Claims []Claim

// Get retrieve existants claims from external DB
func (c *Claims) Get() error {

	return sqlErr
}

// Post a new claim to the external DB
func (c *Claim) Post() error {
	var token []byte
	var database string

	if c.ID == 0 {
		return nil
	}
	token, err := json.Marshal(c.StandardClaims)
	if err != nil {
		fmt.Printf("\n[stack trace] Claim.POST error %s", err.Error())
		return err
	}
	sqlErr = db.QueryRow("SELECT DATABASE()").Scan(&database)
	stmt, sqlErr := claimRequire.Prepare("INSERT INTO claim (idUser, source, claim) values (?,?,?)")
	if sqlErr != nil {
		fmt.Printf("\n[stack trace] Claim.POST error %s", sqlErr.Error())
		return sqlErr
	}
	result, sqlErr := stmt.Exec(c.ID, database, string(token[:]))
	if sqlErr != nil {
		fmt.Printf("\n[stack trace] Claim.POST error %s", sqlErr.Error())
		return sqlErr
	}
	_, err = result.LastInsertId()
	if err != nil {
		fmt.Printf("\n[stack trace] Claim.POST error %s", err.Error())
		return err
	}
	return sqlErr
}

// GetValidClaim check actual claim, if valid, send back, if not create one
func (c *Claim) GetValidClaim(user User, token string) error {
	expireCookie := time.Now().UTC().Add(time.Hour * 1)

	if c.ID == 0 {
		c.ID = user.UID
		c.StandardClaims = jwt.StandardClaims{
			ExpiresAt: expireCookie.Unix(),
			Issuer:    "TIC/WEB_RestFullApi/Castel_a/2019",
		}
		c.Post()
	} else {
		err := json.Unmarshal([]byte(token), &c.StandardClaims)
		if err != nil {
			fmt.Printf("\n[stack trace] < Claim.GetValidClaim() >")
			return err
		}
	}
	return nil
}
