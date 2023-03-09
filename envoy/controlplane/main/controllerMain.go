package main

import (
	"control/plane/controller"
	"control/plane/resource"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func main() {

	go func() {
		fmt.Println("=======================Init Envoy Start=======================")
		//grpcServer.NewServer()

		clusterProp := controller.GetClusterEndpoint()
		snapshot := resource.GenerateSnapshot(clusterProp)
		resource.RefreshSnapshotCache(*snapshot)
		fmt.Println("=======================Init Envoy Finish=======================")
	}()

	app := iris.New()

	app.Get("/add", func(ctx *context.Context) {
		controller.AddClusterEndpoint()
	})
	app.Listen(":80")
}

func initServer() {
	//clusterProp := controller.GetClusterEndpoint()
	//snapshot := resource.GenerateSnapshot(clusterProp)
	//
	//l := controlplane.Logger{Debug: true}
	// Create a cache
	//cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)
	//
	//if err := snapshot.Consistent(); err != nil {
	//	l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	//	os.Exit(1)
	//}
	//
	//// Add the snapshot to the cache
	//if err := cache.SetSnapshot(context.Background(), nodeID, &snapshot); err != nil {
	//	l.Errorf("snapshot error %q for %+v", err, snapshot)
	//	os.Exit(1)
	//}
	//
	//// Run the xDS server
	//ctx := context.Background()
	//cb := &test.Callbacks{Debug: l.Debug}
	//
	//srv := server.NewServer(ctx, cache, cb)
	//grpcServer.RunServer(srv, port)
}
