package resource

import (
	"control/plane/prop"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
)

func GetDefaultClusterEndpoint() *prop.Cluster {
	return &prop.Cluster{Name: ClusterName, SocketAddress: []prop.SocketAddress{{core.SocketAddress_TCP, UpstreamHost, UpstreamPort}}}
}
