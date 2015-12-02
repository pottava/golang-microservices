package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/pottava/golang-micro-services/app-authentication/app/logs"
	"github.com/pottava/golang-micro-services/app-authentication/app/models"
)

const (
	twitterTemporaryKey = "tw-temp"
	twitterSessionKey   = "tw-sess"
)

// CheckTwitterSession returns user information
func CheckTwitterSession(w http.ResponseWriter, r *http.Request) (token models.OAuthAuthorizedToken, tempID string, found bool) {
	if value, err := GetCookie(r, twitterSessionKey); err == nil {
		value = strings.Replace(strings.Replace(strings.Replace(value, ":", "\":\"", -1), "|", ",", -1), ",", "\",\"", -1)
		value = strings.Replace(strings.Replace(strings.Replace(value, "{", "{\"", -1), "}", "\"}", -1), "*", ":", -1)
		if err = json.Unmarshal([]byte(value), &token); err != nil {
			logs.Error.Printf("Error: %v, Value: %v", err, value)
			return
		}
		http.SetCookie(w, SetCookie(twitterTemporaryKey, "", -1))
		return token, "", true
	}
	if value, err := GetCookie(r, twitterTemporaryKey); err == nil {
		return token, value, true
	}
	sessionID, ok := newTwitterSessionID()
	if !ok {
		return token, "", false
	}
	http.SetCookie(w, SetCookie(twitterTemporaryKey, sessionID, 60*60*24))
	return token, "", false
}

func newTwitterSessionID() (string, bool) {
	b := make([]byte, 128)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		logs.Error.Printf("Could not generate new sessionID. Error: %v", err.Error())
		return "", false
	}
	return base64.URLEncoding.EncodeToString(b), true
}
