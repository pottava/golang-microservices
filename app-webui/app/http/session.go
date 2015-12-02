package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Session represents logined user's information
type Session struct {
	UserID string `json:"id"`
}

const (
	sessionKey = "webui-sess"
)

// SetSessionInfo put data to session
func SetSessionInfo(w http.ResponseWriter, id string, duration int) {
	http.SetCookie(w, SetCookie(sessionKey, id, duration))
}

// GetSessionInfo retrives session information
func GetSessionInfo(r *http.Request) *Session {
	session := &Session{}
	if value, err := GetCookie(r, sessionKey); err == nil {
		value = strings.Replace(strings.Replace(strings.Replace(value, ":", "\":\"", -1), "|", ",", -1), ",", "\",\"", -1)
		value = strings.Replace(strings.Replace(strings.Replace(value, "{", "{\"", -1), "}", "\"}", -1), "*", ":", -1)
		if err = json.Unmarshal([]byte(value), &session); err == nil {
			return session
		}
	}
	return &Session{}
}
