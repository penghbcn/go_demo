package controller

import (
	"control/plane/prop"
	"control/plane/resource"
	"errors"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
)

func AddCluster() error {
	return errors.New("当前暂不支持添加cluster")
}

func DeleteClusterEndpoint() {
	//clusterProp, err := GetClusterEndpoint()
	//snapshot := resource.GenerateSnapshot(clusterProp)
	//
	//resource.RefreshSnapshotCache(*snapshot)
}

func AddClusterEndpoint() {
	clusterProp := resource.GetClusterEndpoint(resource.ClusterName)
	clusterProp.SocketAddress = append(clusterProp.SocketAddress, prop.SocketAddress{
		Protocol:  core.SocketAddress_TCP,
		Address:   resource.UpstreamHost,
		PortValue: 8081,
	})
	xds := prop.XDS{Clusters: append(make([]prop.Cluster, 0), *clusterProp)}
	snapshot := resource.GenerateSnapshot(&xds)

	resource.RefreshSnapshotCache(snapshot)
}

func r() {

}
