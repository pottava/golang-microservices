package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pottava/golang-microservices/app-aws/app/logs"
)

const (
	options = "OPTIONS"
	get     = "GET"
	post    = "POST"
	put     = "PUT"
	patch   = "PATCH"
	delete  = "DELETE"
)

// APIStatus represents API's result status
type APIStatus struct {
	success bool
	code    int
	message string
}

// APIResource represents RESTful API Interfaces
type APIResource interface {
	Options(url string, queries url.Values, body io.Reader) (APIStatus, interface{})
	Get(url string, queries url.Values, body io.Reader) (APIStatus, interface{})
	Post(url string, queries url.Values, body io.Reader) (APIStatus, interface{})
	Put(url string, queries url.Values, body io.Reader) (APIStatus, interface{})
	Patch(url string, queries url.Values, body io.Reader) (APIStatus, interface{})
	Delete(url string, queries url.Values, body io.Reader) (APIStatus, interface{})
}

type apiheader struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type apienvelope struct {
	Header   apiheader   `json:"header"`
	Response interface{} `json:"response"`
}

// APIResourceBase is defined for composition
type APIResourceBase struct{}

// Options implements the APIResource Options function
func (APIResourceBase) Options(url string, queries url.Values, body io.Reader) (APIStatus, interface{}) {
	return FailSimple(http.StatusMethodNotAllowed), nil
}

// Get implements the APIResource Get function
func (APIResourceBase) Get(url string, queries url.Values, body io.Reader) (APIStatus, interface{}) {
	return FailSimple(http.StatusMethodNotAllowed), nil
}

// Post implements the APIResource Post function
func (APIResourceBase) Post(url string, queries url.Values, body io.Reader) (APIStatus, interface{}) {
	return FailSimple(http.StatusMethodNotAllowed), nil
}

// Put implements the APIResource Put function
func (APIResourceBase) Put(url string, queries url.Values, body io.Reader) (APIStatus, interface{}) {
	return FailSimple(http.StatusMethodNotAllowed), nil
}

// Patch implements the APIResource Patch function
func (APIResourceBase) Patch(url string, queries url.Values, body io.Reader) (APIStatus, interface{}) {
	return FailSimple(http.StatusMethodNotAllowed), nil
}

// Delete implements the APIResource Delete function
func (APIResourceBase) Delete(url string, queries url.Values, body io.Reader) (APIStatus, interface{}) {
	return FailSimple(http.StatusMethodNotAllowed), nil
}

// APIResourceHandler allows you to implement RESTful APIs easier
func APIResourceHandler(APIResource APIResource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := bytes.NewBuffer(make([]byte, 0))
		reader := io.TeeReader(r.Body, b)

		r.Body = ioutil.NopCloser(b)
		defer r.Body.Close()

		r.ParseForm()

		// Delegate responsibility to the resource
		var status APIStatus
		var data interface{}

		switch r.Method {
		case options:
			status, data = APIResource.Options(r.URL.Path, r.Form, reader)
		case get:
			status, data = APIResource.Get(r.URL.Path, r.Form, reader)
		case post:
			status, data = APIResource.Post(r.URL.Path, r.Form, reader)
		case put:
			status, data = APIResource.Put(r.URL.Path, r.Form, reader)
		case patch:
			status, data = APIResource.Patch(r.URL.Path, r.Form, reader)
		case delete:
			status, data = APIResource.Delete(r.URL.Path, r.Form, reader)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Return API response
		var content []byte
		var e error

		if !status.success {
			content, e = json.Marshal(apienvelope{
				Header: apiheader{Status: "fail", Message: status.message},
			})
		} else {
			content, e = json.Marshal(apienvelope{
				Header:   apiheader{Status: "success"},
				Response: data,
			})
		}
		if e != nil {
			logs.Error.Printf("ERROR: %s %s", "json.Marshal@APIResourceHandler", e.Error())
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.code)
		w.Write(content)
	}
}

// Success means API finished successfully
func Success(code int) APIStatus {
	return APIStatus{success: true, code: code, message: ""}
}

// Fail means API finished unsuccessfully
func Fail(code int, message string) APIStatus {
	return APIStatus{success: false, code: code, message: message}
}

// FailSimple means API finished unsuccessfully
func FailSimple(code int) APIStatus {
	return APIStatus{success: false, code: code, message: strconv.Itoa(code) + " " + http.StatusText(code)}
}
