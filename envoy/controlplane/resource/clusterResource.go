package resource

import (
	"control/plane/prop"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	resourcev3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func GetDefaultClusterEndpoint() *prop.Cluster {
	return &prop.Cluster{Name: ClusterName, SocketAddress: []prop.SocketAddress{{core.SocketAddress_TCP, UpstreamHost, UpstreamPort}}}
}

func GetClusterEndpoint(clusterName string) *prop.Cluster {
	clusterProp := prop.Cluster{Name: clusterName, SocketAddress: make([]prop.SocketAddress, 0)}
	snapshot, err := SnapshotCache.GetSnapshot(NodeId)
	if err != nil {
		return &clusterProp
	}
	resources := snapshot.GetResources(resourcev3.ClusterType)
	resource := resources[clusterName]
	if resource == nil {
		return &clusterProp
	}
	clusterResource := resource.(*cluster.Cluster)
	endpoints := clusterResource.LoadAssignment.Endpoints
	var SocketAddress []prop.SocketAddress
	for _, endpoint := range endpoints {
		for _, lbEndpoint := range endpoint.LbEndpoints {
			addrProp := new(prop.SocketAddress)
			address := lbEndpoint.GetEndpoint().Address.GetSocketAddress()
			addrProp.Protocol = address.Protocol
			addrProp.Address = address.Address
			addrProp.PortValue = address.GetPortValue()
			SocketAddress = append(SocketAddress, *addrProp)
		}
	}
	clusterProp.SocketAddress = SocketAddress
	return &clusterProp
}
