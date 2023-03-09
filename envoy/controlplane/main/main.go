package main

import (
	"control/plane/controller"
	"control/plane/resource"
	grpcServer "control/plane/server"
	"github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/kataras/iris/v12"
	irisContext "github.com/kataras/iris/v12/context"
)

var (
	logger      = log.LoggerFuncs{}
	port   uint = 18000
	nodeID      = "test-id"
	//snapshotCacheCache cache.SnapshotCache
	srv grpcServer.Server
)

func init() {
	go func() {
		resource.InitServer()
	}()
}

func main() {
	app := iris.New()

	app.Get("/add", func(ctx *irisContext.Context) {
		controller.AddClusterEndpoint()
		ctx.Writef("add")
	})

	app.Get("/remove", func(ctx *irisContext.Context) {
		controller.DeleteClusterEndpoint()
		ctx.Writef("remove")
	})

	app.Listen(":80")

	//clusterProp := resource.GetDefaultClusterEndpoint()
	//snapshot := resource.GenerateSnapshot(clusterProp)
	//resource.RefreshSnapshotCache(*snapshot)
}
