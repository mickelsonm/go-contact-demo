package controllers

import (
	"github.com/ninnemana/web"
)

func Index(ctx *web.Context, args ...string) {

	tmpl := web.NewTemplate(ctx.ResponseWriter)
	tmpl.ParseFile("templates/index.html", false)

	tmpl.Display(ctx.ResponseWriter)
}
