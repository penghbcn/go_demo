package resource

import (
	"context"
	"control/plane/prop"
	grpcServer "control/plane/server"
	"control/plane/version"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"os"
)

var (
	logger      = log.LoggerFuncs{}
	port   uint = 18000
	nodeID      = "test-id"

	snapshotCacheCache cache.SnapshotCache
	srv                grpcServer.Server
)

func InitServer() {
	// Create a cache
	snapshotCacheCache = cache.NewSnapshotCache(true, cache.IDHash{}, logger)
	clusterProp := GetDefaultClusterEndpoint()
	snapshot := GenerateSnapshot(clusterProp)
	// Add the snapshot to the cache
	if err := snapshotCacheCache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		logger.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	// Run the xDS server
	ctx := context.Background()
	cb := &test.Callbacks{Debug: true}

	srv = *grpcServer.NewServer(ctx, snapshotCacheCache, cb)
	srv.Run(port)
}

func RefreshSnapshotCache(snapshot cache.Snapshot) {

	ctx := context.Background()

	err := snapshotCacheCache.SetSnapshot(ctx, nodeID, &snapshot)
	if err != nil {
		logger.Errorf("", err)
	}
	//srv
	//// Create a cache
	//cache := cache.NewSnapshotCache(false, cache.IDHash{}, logger)
	//
	//if err := snapshot.Consistent(); err != nil {
	//	logger.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
	//	os.Exit(1)
	//}
	//logger.Debugf("will serve snapshot %+v", snapshot)
	//
	//// Add the snapshot to the cache
	//if err := main.snapshotCacheCache.SetSnapshot(context.Background(), nodeID, &snapshot); err != nil {
	//	logger.Errorf("snapshot error %q for %+v", err, snapshot)
	//	os.Exit(1)
	//}
	//
	// Run the xDS server
	//ctx := context.Background()
	//cb := &test.Callbacks{Debug: true}
	//
	//srv := server.NewServer(ctx, cache, cb)
	//grpcServer.RunServer(srv, port)
}

func GenerateSnapshot(clusterProp *prop.Cluster) *cache.Snapshot {
	snap, _ := cache.NewSnapshot(version.GetNewVersion(),
		map[resource.Type][]types.Resource{
			resource.ClusterType:  {makeCluster(clusterProp)},
			resource.ListenerType: {makeHTTPListener(ListenerName, RouteName)},
		},
	)
	return snap
}
