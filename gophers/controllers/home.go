package controllers

import (
	"gophers/plate"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	server := plate.NewServer()

	tmpl, _ := server.Template(w)
	tmpl.Template = "templates/index.html"

	tmpl.DisplayTemplate()
}

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	plate.Serve404(w, "")
}
