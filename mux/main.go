package main

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

//go:embed template/*.html
var tpl embed.FS

func main() {
	ct, err := tpl.ReadFile("template/index.html")
	fmt.Printf("%s - %v\n", ct, err)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			slog.Warn(fmt.Sprintf("parse tpl err: %s", err))
		}
		err = tpl.Execute(w, nil)
		if err != nil {
			slog.Warn(fmt.Sprintf("output tpl err: %s", err))
		}
	})

	proxy, _ := url.Parse("http://127.0.0.1:10809")
	http.DefaultTransport.(*http.Transport).Proxy = http.ProxyURL(proxy)
	//parseMenu()

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("run server error: %s", err)
	}
}
