package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/kurrik/oauth1a"
	"github.com/pottava/golang-microservices/app-authentication/app/config"
	util "github.com/pottava/golang-microservices/app-authentication/app/http"
	"github.com/pottava/golang-microservices/app-authentication/app/logs"
	"github.com/pottava/golang-microservices/app-authentication/app/models"
)

var (
	sessions map[string]*oauth1a.UserConfig
	twitter  *oauth1a.Service
)

const (
	twitterTemporaryKey = "tw-temp"
	twitterSessionKey   = "tw-sess"
)

func init() {
	sessions = map[string]*oauth1a.UserConfig{}
	cfg := config.NewConfig()
	twitter = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    cfg.TwitterKey,
			ConsumerSecret: cfg.TwitterSecret,
			CallbackURL:    cfg.TwitterCallback,
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}
	re := regexp.MustCompile("https*:")

	/**
	 * Twitter OAuth
	 */
	// out -> Location: /twitter/signin
	//        Set-Cookie: tw-temp=X
	http.Handle("/twitter", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		util.CheckTwitterSession(w, r)
		http.Redirect(w, r, "/twitter/signin", http.StatusFound)
	}))

	// in  -> Cookie: tw-temp=X
	// out -> Location: https://api.twitter.com/oauth/authorize?oauth_token=X
	//        Set-Cookie: tw-sess=
	http.Handle("/twitter/signin", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, util.SetCookie(twitterSessionKey, "", -1))
		_, sessionID, ok := util.CheckTwitterSession(w, r)
		if !ok {
			logs.Error.Print("Could not find temporary session ID.")
			http.Error(w, "Problem generating new session", http.StatusInternalServerError)
			return
		}
		session := &oauth1a.UserConfig{}
		if err := session.GetRequestToken(twitter, new(http.Client)); err != nil {
			logs.Error.Printf("Could not get request token: %v", err)
			http.Error(w, fmt.Sprintf("Problem getting the request token: %v", err), http.StatusInternalServerError)
			return
		}
		url, err := session.GetAuthorizeURL(twitter)
		if err != nil {
			logs.Error.Printf("Could not get authorization URL: %v", err)
			http.Error(w, "Problem getting the authorization URL", http.StatusInternalServerError)
			return
		}
		sessions[sessionID] = session
		http.Redirect(w, r, url, http.StatusFound)
	}))

	// in  -> Cookie: tw-temp=X
	//        /twitter/callback?oauth_token=X&oauth_verifier=Y
	//        /twitter/callback?denied=Z
	// out -> Location: /
	//        Set-Cookie: tw-sess=Y
	http.Handle("/twitter/callback", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		_, found := util.RequestGetParam(r, "denied")
		if found {
			http.Redirect(w, r, "http://192.168.99.100/", http.StatusFound)
			return
		}
		authorizedToken, sessionID, ok := util.CheckTwitterSession(w, r)
		if !ok {
			logs.Error.Print("Could not find any authorized token.")
			http.Error(w, "Problem generating new session", http.StatusInternalServerError)
			return
		}
		if authorizedToken.AccessTokenKey != "" {
			http.Redirect(w, r, "http://192.168.99.100/", http.StatusFound)
			return
		}
		session, ok := sessions[sessionID]
		if !ok {
			logs.Error.Print("Could not find user config in sesions storage.")
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}
		token, verifier, err := session.ParseAuthorize(r, twitter)
		if err != nil {
			logs.Error.Printf("Could not parse authorization: %v", err)
			http.Error(w, "Problem parsing authorization", http.StatusInternalServerError)
			return
		}
		if err = session.GetAccessToken(token, verifier, twitter, new(http.Client)); err != nil {
			logs.Error.Printf("Error getting access token: %v", err)
			http.Error(w, "Problem getting an access token", http.StatusInternalServerError)
			return
		}
		delete(sessions, sessionID)

		// Update user information
		userID := "tw/" + session.AccessValues.Get("user_id")
		user, found := models.GetUser(userID)
		if !found {
			user = &models.User{}
			user.ID = userID
		}
		user.Name = session.AccessValues.Get("screen_name")
		models.SaveUser(user)

		// Set OAuth session
		authorized := models.OAuthAuthorizedToken{
			UserID:            userID,
			ScreenName:        session.AccessValues.Get("screen_name"),
			AccessTokenKey:    session.AccessTokenKey,
			AccessTokenSecret: session.AccessTokenSecret,
		}
		authorized.ScreenImage = re.ReplaceAllLiteralString(twitterImage(authorized), "")

		bytes, _ := json.Marshal(authorized)
		formatted := strings.Replace(string(bytes), ",", "|", -1)
		http.SetCookie(w, util.SetCookie(twitterSessionKey, formatted, 60*60*24))
		http.Redirect(w, r, "http://192.168.99.100/auth/login?id="+authorized.UserID, http.StatusFound)
	}))

	// out -> Location: /
	//        Set-Cookie: tw-sess=
	http.Handle("/twitter/logout", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, util.SetCookie(twitterSessionKey, "", -1))
		http.Redirect(w, r, "http://192.168.99.100/auth/logout", http.StatusFound)
	}))
}

func twitterImage(sess models.OAuthAuthorizedToken) string {
	query := url.Values{}
	query.Set("screen_name", sess.ScreenName)
	url := fmt.Sprintf("https://api.twitter.com/1.1/users/show.json?%v", query.Encode())

	req, e1 := http.NewRequest("GET", url, nil)
	if e1 != nil {
		logs.Error.Printf("Could not get new request: %v", e1)
		return ""
	}
	twitter.Sign(req, oauth1a.NewAuthorizedConfig(sess.AccessTokenKey, sess.AccessTokenSecret))

	client := &http.Client{}
	res, e2 := client.Do(req)
	if e2 != nil {
		logs.Error.Printf("Could not send HTTP request: %v", e2)
		return ""
	}
	defer res.Body.Close()

	body, e3 := ioutil.ReadAll(res.Body)
	if e3 != nil {
		logs.Error.Printf("Could not read HTTP response: %v", e3)
		return ""
	}
	type user struct {
		ProfileImageURL string `json:"profile_image_url"`
	}
	var u user
	if err := json.Unmarshal(body, &u); err != nil {
		logs.Error.Printf("Error: %v", err)
		return ""
	}
	return u.ProfileImageURL
}
