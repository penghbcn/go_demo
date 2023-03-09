package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func main() {
	fmt.Println("=======================")

	app := iris.New()

	app.Get("/{path}", func(ctx *context.Context) {
		path := ctx.Params().Get("path")
		ctx.Writef(path + "  8081")
	})
	app.Listen(":8081")
}
