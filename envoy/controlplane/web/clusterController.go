package web

import (
	"control/plane/controller"
	"control/plane/resource"
	irisContext "github.com/kataras/iris/v12/context"
)

func clusterController() {
	app.Get("/get", func(ctx *irisContext.Context) {
		endpoint := resource.GetClusterEndpoint(resource.ClusterName)
		ctx.JSON(endpoint)

	})

	app.Get("/add", func(ctx *irisContext.Context) {
		controller.AddClusterEndpoint()
		ctx.Writef("add")
	})

	app.Get("/remove", func(ctx *irisContext.Context) {
		controller.DeleteClusterEndpoint()
		ctx.Writef("remove")
	})
}
