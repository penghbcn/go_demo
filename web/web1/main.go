package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func main() {
	app := iris.New()

	app.Get("/{path}", func(ctx *context.Context) {
		path := ctx.Params().Get("path")
		ctx.Writef(path + "  8080")
	})
	println("======================= Web Server Listening On 8080 =======================")
	app.Listen(":8080")
}
