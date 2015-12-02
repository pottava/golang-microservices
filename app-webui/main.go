package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/pottava/golang-micro-services/app-webui/app/config"
	"github.com/pottava/golang-micro-services/app-webui/app/controllers"
	util "github.com/pottava/golang-micro-services/app-webui/app/http"
	"github.com/pottava/golang-micro-services/app-webui/app/logs"
)

func main() {
	cfg := config.NewConfig()

	http.Handle("/", index(cfg))
	http.Handle("/assets/", assets(cfg))

	logs.Debug.Print("[config] " + cfg.String())
	logs.Info.Printf("[service] listening on port %v", cfg.Port)
	logs.Fatal.Print(http.ListenAndServe(":"+fmt.Sprint(cfg.Port), nil))
}

func index(cfg *config.Config) http.Handler {
	return util.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		util.RenderHTML(w, []string{"app/index.tmpl"}, controllers.CommonParams(w, r), nil)
	})
}

func assets(cfg *config.Config) http.Handler {
	fs := http.FileServer(http.Dir(path.Join(cfg.StaticFilePath, "assets")))
	return util.AssetsChain(http.StripPrefix("/assets/", fs).ServeHTTP)
}
