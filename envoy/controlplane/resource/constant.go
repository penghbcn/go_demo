package resource

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/log"
)

const (
	NodeId       = "test-id"
	ClusterName  = "default_proxy_cluster"
	RouteName    = "default_route"
	ListenerName = "default_listener_0"
	ListenerPort = 10000
	UpstreamHost = "192.168.12.91"
	UpstreamPort = 8080
)

var (
	logger = log.LoggerFuncs{}
	nodeID = "test-id"

	SnapshotCache cache.SnapshotCache
)
