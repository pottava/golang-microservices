package models

// OAuthAuthorizedToken represents authorized person
type OAuthAuthorizedToken struct {
	UserID            string `json:"id"`
	ScreenName        string `json:"name"`
	ScreenImage       string `json:"img"`
	AccessTokenKey    string `json:"token"`
	AccessTokenSecret string `json:"secret"`
}
