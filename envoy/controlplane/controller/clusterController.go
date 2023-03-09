package controller

import (
	"control/plane/prop"
	"control/plane/resource"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
)

const (
	ClusterName  = "example_proxy_cluster"
	RouteName    = "local_route"
	ListenerName = "listener_0"
	ListenerPort = 10000
	UpstreamHost = "192.168.12.91"
	UpstreamPort = 8080
)

func AddCluster() {

}

func DeleteClusterEndpoint() {
	clusterProp := GetClusterEndpoint()
	snapshot := resource.GenerateSnapshot(clusterProp)

	resource.RefreshSnapshotCache(*snapshot)
}

func AddClusterEndpoint() {
	clusterProp := GetClusterEndpoint()
	clusterProp.SocketAddress = append(clusterProp.SocketAddress, prop.SocketAddress{
		Protocol:  core.SocketAddress_TCP,
		Address:   UpstreamHost,
		PortValue: 8081,
	})
	snapshot := resource.GenerateSnapshot(clusterProp)

	resource.RefreshSnapshotCache(*snapshot)
}

func GetClusterEndpoint() *prop.Cluster {
	return &prop.Cluster{Name: ClusterName, SocketAddress: []prop.SocketAddress{{core.SocketAddress_TCP, UpstreamHost, UpstreamPort}}}
}
