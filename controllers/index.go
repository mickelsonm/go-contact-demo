package controllers

import (
	"../helpers/database"
	"github.com/ninnemana/web"
	"net/http"
	"time"
)

type Contact struct {
	Id                           int
	Fname, Lname, Email, Message string
	Created                      time.Time
}

func GetMessage(ctx *web.Context, id string) string {
	if id != "" {
		sel, err := database.GetStatement("GetMessageById")

		if err != nil {
			ctx.Abort(http.StatusInternalServerError, err.Error())
			return ""
		}

		sel.Bind(id)

		rows, res, err := sel.Exec()

		if err != nil {
			ctx.Abort(http.StatusInternalServerError, err.Error())
			return ""
		}

		message := res.Map("message")

		for _, row := range rows {
			return row.Str(message)
		}
	}

	return ""
}

func Add(ctx *web.Context) {
	var c Contact

	c.Fname = ctx.Request.FormValue("fname")
	c.Lname = ctx.Request.FormValue("lname")
	c.Email = ctx.Request.FormValue("email")
	c.Message = ctx.Request.FormValue("message")

	ins, err := database.GetStatement("InsertContact")

	if err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
		return
	}

	ins.Bind(c.Fname, c.Lname, c.Email, c.Message, time.Now())

	_, _, err = ins.Exec()

	if err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

func Index(ctx *web.Context, args ...string) {

	var contacts []Contact

	sel, err := database.GetStatement("GetAllContacts")

	if err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
	}

	rows, res, err := sel.Exec()

	id := res.Map("id")
	fname := res.Map("fname")
	lname := res.Map("lname")
	email := res.Map("email")
	message := res.Map("message")
	created := res.Map("created")

	for _, row := range rows {
		c := Contact{
			Id:      row.Int(id),
			Fname:   row.Str(fname),
			Lname:   row.Str(lname),
			Email:   row.Str(email),
			Message: row.Str(message),
			Created: row.ForceLocaltime(created),
		}

		contacts = append(contacts, c)
	}

	tmpl := web.NewTemplate(ctx.ResponseWriter)
	tmpl.ParseFile("templates/index.html", false)
	tmpl.Bag["Contacts"] = contacts

	tmpl.Display(ctx.ResponseWriter)
}
