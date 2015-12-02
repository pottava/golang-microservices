// Package http provides a convenient way to impliment http servers
package http

import (
	"compress/gzip"
	"compress/zlib"
	"html"
	"io"
	"net/http"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/justinas/alice"
	"github.com/pottava/golang-microservices/app-webui/app/config"
	"github.com/pottava/golang-microservices/app-webui/app/logs"
	"github.com/pottava/golang-microservices/app-webui/app/misc"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

// RequestGetParam retrives a request parameter
func RequestGetParam(r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	return value, (len(value) != 0)
}

// SetCookie set a cookie
func SetCookie(key, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: maxAge,
		Secure: cfg.SecuredCookie,
		Path:   "/",
	}
}

// GetCookie retrives a cookie value
func GetCookie(r *http.Request, key string) (id string, err error) {
	var c *http.Cookie
	if c, err = r.Cookie(key); err != nil {
		return
	}
	id = c.Value
	return
}

// Chain enables middleware chaining
func Chain(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return chain(true, true, true, f)
}

// AssetsChain enables middleware chaining
func AssetsChain(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return chain(false, true, false, f)
}

// RenderText write data as a simple text
func RenderText(w http.ResponseWriter, data string, err error) {
	if isInvalid(w, err, "@RenderText") {
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(html.EscapeString(data)))
}

// RenderHTML write data as a HTML text with template
func RenderHTML(w http.ResponseWriter, templatePath []string, data interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relatives := append([]string{"base.tmpl"}, templatePath...)
	templates := make([]string, len(relatives))
	for idx, template := range relatives {
		templates[idx] = path.Join(cfg.StaticFilePath, "views", template)
	}

	tmpl, err := template.ParseFiles(templates...)
	if isInvalid(w, err, "@RenderHTML") {
		return
	}
	if err := tmpl.Execute(w, struct {
		AppName        string
		StaticFileHost string
		Data           interface{}
	}{cfg.Name, cfg.StaticFileHost, data}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error.Printf("ERROR: @RenderHTML %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html")
}

func isInvalid(w http.ResponseWriter, err error, caption string) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logs.Error.Printf("ERROR: %s %s", caption, err.Error())
		return true
	}
	return false
}

type customResponseWriter struct {
	io.Writer
	http.ResponseWriter
	status int
}

func (r *customResponseWriter) Write(b []byte) (int, error) {
	if r.Header().Get("Content-Type") == "" {
		r.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return r.Writer.Write(b)
}

func (r *customResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}

func chain(log, cors, validate bool, f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return alice.New(timeout).Then(http.HandlerFunc(custom(log, cors, validate, f)))
}

func custom(log, cors, validate bool, f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr
		if ip, found := header(r, "X-Forwarded-For"); found {
			addr = ip
		}
		// compress settings
		ioWriter := w.(io.Writer)
		for _, val := range misc.ParseCsvLine(r.Header.Get("Accept-Encoding")) {
			if val == "gzip" {
				w.Header().Set("Content-Encoding", "gzip")
				g := gzip.NewWriter(w)
				defer g.Close()
				ioWriter = g
				break
			}
			if val == "deflate" {
				w.Header().Set("Content-Encoding", "deflate")
				z := zlib.NewWriter(w)
				defer z.Close()
				ioWriter = z
				break
			}
		}
		writer := &customResponseWriter{Writer: ioWriter, ResponseWriter: w, status: http.StatusOK}

		// route to the controllers
		f(writer, r)

		// access log
		if log && cfg.AccessLog {
			logs.Info.Printf("%s %s %s %s", addr, strconv.Itoa(writer.status), r.Method, r.URL)
		}
	}
}

func header(r *http.Request, key string) (string, bool) {
	if r.Header == nil {
		return "", false
	}
	if candidate := r.Header[key]; !misc.ZeroOrNil(candidate) {
		return candidate[0], true
	}
	return "", false
}

func timeout(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 300*time.Second, "timed out")
}
