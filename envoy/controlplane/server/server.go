package grpcServer

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

type Server struct {
	xdsServer server.Server
}

func NewServer(ctx context.Context, cache cache.SnapshotCache, cb *test.Callbacks) *Server {
	srv := server.NewServer(ctx, cache, cb)
	return &Server{srv}
}

func (s *Server) registerServer(grpcServer *grpc.Server) {
	// register services
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(grpcServer, s.xdsServer)
	endpointv3.RegisterEndpointDiscoveryServiceServer(grpcServer, s.xdsServer)
	clusterv3.RegisterClusterDiscoveryServiceServer(grpcServer, s.xdsServer)
	routev3.RegisterRouteDiscoveryServiceServer(grpcServer, s.xdsServer)
	listenerv3.RegisterListenerDiscoveryServiceServer(grpcServer, s.xdsServer)
	secretv3.RegisterSecretDiscoveryServiceServer(grpcServer, s.xdsServer)
	runtimev3.RegisterRuntimeDiscoveryServiceServer(grpcServer, s.xdsServer)
}

func (s *Server) Run(port uint) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems. Keepalive timeouts based on connection_keepalive parameter https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/examples#dynamic
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s.registerServer(grpcServer)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}
