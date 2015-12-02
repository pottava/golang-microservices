package controllers

import (
	"io"
	"net/http"
	"net/url"

	util "github.com/pottava/golang-microservices/app-dbio/app/http"
	"github.com/pottava/golang-microservices/app-dbio/app/logs"
	"github.com/pottava/golang-microservices/app-dbio/app/misc"
	"github.com/pottava/golang-microservices/app-dbio/app/models"
)

func init() {
	http.Handle("/users/", util.Chain(util.APIResourceHandler(users{})))
}

type users struct {
	util.APIResourceBase
}

func (c users) Get(url string, queries url.Values, body io.Reader) (util.APIStatus, interface{}) {
	// retrive a specified user
	if id := url[len("/users/"):]; len(id) != 0 {
		user, found := models.GetUser(id)
		if !found {
			return util.Success(http.StatusOK), models.User{}
		}
		return util.Success(http.StatusOK), user
	}
	// list users
	users, _, err := models.GetUsers()
	if err != nil {
		return util.Fail(http.StatusInternalServerError, err.Error()), nil
	}
	return util.Success(http.StatusOK), users
}

func (c users) Post(url string, queries url.Values, body io.Reader) (util.APIStatus, interface{}) {
	user := &models.User{}
	if err := misc.ReadMBJSON(body, user, 100); err != nil {
		logs.Error.Printf("Could not decode response body as a json. Error: %v", err)
		return util.Fail(http.StatusInternalServerError, err.Error()), nil
	}
	if err := user.Persist(); err != nil {
		return util.Fail(http.StatusInternalServerError, err.Error()), nil
	}
	return util.Success(http.StatusOK), user
}
