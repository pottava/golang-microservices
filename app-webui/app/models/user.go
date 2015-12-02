package models

// User represents user's user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type daoUser struct {
	Header   APIHeader `json:"header"`
	Response User      `json:"response"`
}

// UserParameters represents user information parameters
type UserParameters struct {
	ID string
}

// PageParameters represents html parameters
type PageParameters struct {
	User *UserParameters
	Mode string
}

// GetUser retrives a specified user from DynamoDB
//  @param  id string
//  @return user models.User
func GetUser(userID string) (user *User, found bool) {
	res := &daoUser{}
	if err := db("GET", "/users/"+userID, "", res); err == nil {
		if res.Header.Status == "success" && res.Response.ID == userID {
			return &res.Response, true
		}
	}
	return &User{}, false
}
