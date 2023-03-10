package web

import (
	"github.com/kataras/iris/v12"
)

var app *iris.Application

func RunWebServer() {
	app = iris.New()

	app.Listen(":80")
}
