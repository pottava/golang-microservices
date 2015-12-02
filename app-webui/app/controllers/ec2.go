// Package controllers implements functions to route user requests
package controllers

import (
	"io"
	"net/http"
	"net/url"

	util "github.com/pottava/golang-microservices/app-webui/app/http"
	"github.com/pottava/golang-microservices/app-webui/app/misc"
	"github.com/pottava/golang-microservices/app-webui/app/models"
)

func init() {
	http.Handle("/ec2/instances/", util.Chain(util.APIResourceHandler(ec2Instances{})))
}

type ec2Instances struct {
	util.APIResourceBase
}

func (c ec2Instances) Get(session *util.Session, url string, queries url.Values, body io.Reader) (util.APIStatus, interface{}) {
	if misc.ZeroOrNil(session) || misc.ZeroOrNil(session.UserID) {
		return util.FailSimple(http.StatusUnauthorized), nil
	}
	instances, _ := models.GetEC2Instances()
	return util.Success(http.StatusOK), struct {
		EC2Instances []*models.EC2Instance `json:"instances"`
		Count        int                   `json:"count"`
	}{
		EC2Instances: instances,
		Count:        len(instances),
	}
}
