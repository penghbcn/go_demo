package prop

import "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

type Cluster struct {
	Name          string
	SocketAddress []SocketAddress
}

type SocketAddress struct {
	Protocol  corev3.SocketAddress_Protocol
	Address   string
	PortValue uint32
}
