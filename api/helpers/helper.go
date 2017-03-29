package helpers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ParseAuthorization decode authorization header content and send it back into the response as json object.
func ParseAuthorization(r *http.Request) ([]string, error) {
	var parseErr ParseErr

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		parseErr.Message = "Header Authorization do not permit to identify BasicHttpAuthentication."
		fmt.Printf("\n[basic security] Header Authorization do not permit to identify BasicHttpAuthentication.\n")
		fmt.Println(auth)
		return nil, parseErr
	}
	decodedCredentials, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		fmt.Printf("\n[basic security] base64 encoded 'Authorization' content not properly decoded.")
		fmt.Printf("\n %s", auth[1])
		return nil, err
	}
	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) == 2 {
		credentials[1], err = ToSha1(credentials[1])
	}
	return credentials, err
}

// ToSha1 encode string to sha1 and return the encoded string
func ToSha1(toEncode string) (string, error) {
	if toEncode == "" {
		err := errors.New("\nNothing to encode in sha1")
		return "", err
	}
	tmp := sha1.New()
	_, err := io.WriteString(tmp, toEncode)
	if err != nil {
		return "", err
	}
	toEncode = hex.EncodeToString(tmp.Sum(nil))
	return toEncode, nil
}
