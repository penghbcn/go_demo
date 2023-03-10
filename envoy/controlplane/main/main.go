package main

import (
	"context"
	"control/plane/controller"
	"control/plane/prop"
	"control/plane/resource"
	grpcServer "control/plane/server"
	"control/plane/web"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"github.com/kataras/iris/v12"
	irisContext "github.com/kataras/iris/v12/context"
	"os"
)

var (
	logger      = log.LoggerFuncs{}
	port   uint = 18000
	nodeID      = "test-id"
)

func main() {

	go func() {
		web.RunWebServer()
	}()

	runXDSServer(nodeID, port)
}

func runWebServer() {
	app := iris.New()

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

	app.Listen(":80")
}

func runXDSServer(nodeId string, port uint) {
	// 创建一个缓存
	snapshotCache := cache.NewSnapshotCache(true, cache.IDHash{}, logger)
	// 保存到成员变量,用于后续动态更改
	resource.SnapshotCache = snapshotCache
	// 生成snapshot
	clusterProp := resource.GetDefaultClusterEndpoint()

	xds := prop.XDS{Clusters: append(make([]prop.Cluster, 0), *clusterProp)}
	snapshot := resource.GenerateSnapshot(&xds)
	// 校验
	if err := snapshot.Consistent(); err != nil {
		logger.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}
	// 将snapshot保存到缓存
	if err := snapshotCache.SetSnapshot(context.Background(), nodeId, snapshot); err != nil {
		logger.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	// Run the xDS server
	ctx := context.Background()
	cb := &test.Callbacks{Debug: true}
	srv := *grpcServer.NewServer(ctx, snapshotCache, cb)
	srv.Run(port)
}
