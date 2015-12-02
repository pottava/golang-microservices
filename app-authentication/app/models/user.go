package models

import (
	"encoding/json"
	"errors"
)

// User represents user's user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type daoUser struct {
	Header   APIHeader `json:"header"`
	Response User      `json:"response"`
}

// GetUser retrives a specified user from DynamoDB
//  @param  id string
//  @return user models.User
func GetUser(userID string) (user *User, found bool) {
	res := &daoUser{}
	if db("GET", "/users/"+userID, "", res) == nil {
		if res.Header.Status == "success" && res.Response.ID == userID {
			return &res.Response, true
		}
	}
	return &User{}, false
}

// SaveUser persist its state
//  @param user models.User
func SaveUser(user *User) error {
	req, err := json.Marshal(user)
	if err != nil {
		return err
	}
	res := &daoUser{}
	err = db("POST", "/users/", string(req), res)
	if err != nil {
		return err
	}
	if res.Header.Status == "success" && res.Response.ID == user.ID {
		return nil
	}
	return errors.New(res.Header.Message)
}
