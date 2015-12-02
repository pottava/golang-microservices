package main

import (
	"fmt"
	"net/http"

	"github.com/pottava/golang-micro-services/app-authentication/app/config"
	_ "github.com/pottava/golang-micro-services/app-authentication/app/controllers"
	"github.com/pottava/golang-micro-services/app-authentication/app/logs"
)

func main() {
	cfg := config.NewConfig()
	logs.Debug.Print("[config] " + cfg.String())
	logs.Info.Printf("[service] listening on port %v", cfg.Port)
	logs.Fatal.Print(http.ListenAndServe(":"+fmt.Sprint(cfg.Port), nil))
}
