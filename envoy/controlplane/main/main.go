package main

//
//import (
//	"context"
//	plane "control/plane"
//	"control/plane/controller"
//	"control/plane/resource"
//	"control/plane/server"
//
//	"flag"
//	"os"
//
//	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
//	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
//	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
//)
//
//var (
//	l      plane.Logger
//	port   uint
//	nodeID string
//)
//
//func init() {
//	l = plane.Logger{}
//
//	flag.BoolVar(&l.Debug, "debug", true, "Enable xDS server debug logging")
//
//	// The port that this xDS server listens on
//	flag.UintVar(&port, "port", 18000, "xDS management server port")
//
//	// Tell Envoy to use this Node ID
//	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")
//}
//
//func main() {
//	clusterProp := controller.GetClusterEndpoint()
//	snapshot := resource.GenerateSnapshot(clusterProp)
//	RefreshSnapshotCache(snapshot)
//}
//
//func RefreshSnapshotCache(snapshot *cache.Snapshot) {
//	flag.Parse()
//
//	// Create a cache
//	cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)
//
//	if err := snapshot.Consistent(); err != nil {
//		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
//		os.Exit(1)
//	}
//	l.Debugf("will serve snapshot %+v", snapshot)
//
//	// Add the snapshot to the cache
//	if err := cache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
//		l.Errorf("snapshot error %q for %+v", err, snapshot)
//		os.Exit(1)
//	}
//
//	// Run the xDS server
//	ctx := context.Background()
//	cb := &test.Callbacks{Debug: l.Debug}
//	srv := server.NewServer(ctx, cache, cb)
//	grpcServer.RunServer(srv, port)
//}
