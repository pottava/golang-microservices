// Package controllers implements functions to route user requests
package controllers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pottava/golang-microservices/app-aws/app/aws"
	util "github.com/pottava/golang-microservices/app-aws/app/http"
)

func init() {
	http.Handle("/ec2/instances/", util.Chain(util.APIResourceHandler(ec2Instances{})))
}

type ec2Instances struct {
	util.APIResourceBase
}

func (c ec2Instances) Get(url string, queries url.Values, body io.Reader) (util.APIStatus, interface{}) {
	// retrive a specified instance
	if id := url[len("/ec2/instances/"):]; len(id) != 0 {
		instance, err := aws.Ec2Instance(id)
		if err != nil {
			return util.Fail(http.StatusInternalServerError, err.Error()), nil
		}
		return util.Success(http.StatusOK), instance
	}
	// list instances
	instances, err := aws.Ec2Instances()
	if err != nil {
		return util.Fail(http.StatusInternalServerError, err.Error()), nil
	}
	return util.Success(http.StatusOK), instances
}
