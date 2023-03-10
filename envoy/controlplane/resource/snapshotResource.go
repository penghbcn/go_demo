package resource

import (
	"context"
	"control/plane/prop"
	"control/plane/version"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func RefreshSnapshotCache(snapshot *cache.Snapshot) {
	ctx := context.Background()
	// 将snapshot存到cache中,会自动同步到envoy
	err := SnapshotCache.SetSnapshot(ctx, nodeID, snapshot)
	if err != nil {
		logger.Errorf("%s", err)
	}
}

func GenerateSnapshot(xDSProp *prop.XDS) *cache.Snapshot {
	snap, _ := cache.NewSnapshot(version.GetNewVersion(),
		map[resource.Type][]types.Resource{
			resource.ClusterType:  makeClusters(xDSProp.Clusters),
			resource.ListenerType: {makeHTTPListener(ListenerName, RouteName)},
		},
	)

	return snap
}
