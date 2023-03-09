package resource

import (
	"control/plane/prop"
	"control/plane/version"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func GenerateBaseSnapshot() *cache.Snapshot {
	var clusterProp = prop.Cluster{Name: ClusterName, SocketAddress: []prop.SocketAddress{{core.SocketAddress_TCP, UpstreamHost, UpstreamPort}}}

	snap, _ := cache.NewSnapshot(version.GetNewVersion(),
		map[resource.Type][]types.Resource{
			resource.ClusterType:  {makeCluster(&clusterProp)},
			resource.ListenerType: {makeHTTPListener(ListenerName, RouteName)},
		},
	)
	return snap
}
