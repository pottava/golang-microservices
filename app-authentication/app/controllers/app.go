package controllers

import (
	"io"
	"net/http"
	"net/url"

	util "github.com/pottava/golang-microservices/app-authentication/app/http"
	"github.com/pottava/golang-microservices/app-authentication/app/misc"
	"github.com/pottava/golang-microservices/app-authentication/app/models"
)

func init() {
	http.Handle("/authenticated", util.Chain(util.APIResourceHandler(auth{})))
}

type auth struct {
	util.APIResourceBase
}

func (c auth) Get(session *models.OAuthAuthorizedToken, url string, queries url.Values, body io.Reader) (util.APIStatus, interface{}) {
	user := &struct {
		Name        string `json:"name"`
		Image       string `json:"img, omitempty"`
		TokenKey    string `json:"token, omitempty"`
		TokenSecret string `json:"secret, omitempty"`
	}{}
	if !misc.ZeroOrNil(session) {
		user.Name = session.ScreenName
		user.Image = session.ScreenImage
		user.TokenKey = session.AccessTokenKey
		user.TokenSecret = session.AccessTokenSecret
	}
	return util.Success(http.StatusOK), user
}
