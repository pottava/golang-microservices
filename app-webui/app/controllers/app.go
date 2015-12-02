package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pottava/golang-microservices/app-webui/app/config"
	util "github.com/pottava/golang-microservices/app-webui/app/http"
	"github.com/pottava/golang-microservices/app-webui/app/models"
)

func init() {

	http.Handle("/auth/login", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		id, found := util.RequestGetParam(r, "id")
		if found {
			bytes, _ := json.Marshal(util.Session{UserID: id})
			formatted := strings.Replace(string(bytes), ",", "|", -1)
			util.SetSessionInfo(w, formatted, 60*60*24)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}))

	http.Handle("/auth/logout", util.Chain(func(w http.ResponseWriter, r *http.Request) {
		util.SetSessionInfo(w, "", -1)
		http.Redirect(w, r, "/", http.StatusFound)
	}))
}

// CommonParams returns html common parameters
func CommonParams(w http.ResponseWriter, r *http.Request) *models.PageParameters {
	session := util.GetSessionInfo(r)
	cfg := config.NewConfig()

	return &models.PageParameters{
		User: &models.UserParameters{
			ID: session.UserID,
		},
		Mode: cfg.Mode,
	}
}
