// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package resource

import (
	"control/plane/prop"
	"github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/local_ratelimit/v3"
	"github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
)

const (
	ClusterName  = "example_proxy_cluster"
	RouteName    = "local_route"
	ListenerName = "listener_0"
	ListenerPort = 10000
	UpstreamHost = "192.168.12.91"
	UpstreamPort = 8080
)

func makeCluster(clusterProp *prop.Cluster) *cluster.Cluster {

	return &cluster.Cluster{
		Name:                 clusterProp.Name,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STATIC},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeLoadAssignment(clusterProp),
		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
	}
}

func makeLbEndpoint(socketAddress []prop.SocketAddress) []*endpoint.LbEndpoint {
	var lbEndpoints []*endpoint.LbEndpoint
	for _, addr := range socketAddress {
		lbEndpoint := &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Protocol: core.SocketAddress_TCP,
								Address:  addr.Address,
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: addr.PortValue,
								},
							},
						},
					},
				},
			},
		}
		lbEndpoints = append(lbEndpoints, lbEndpoint)
	}
	return lbEndpoints
}

func makeLoadAssignment(clusterProp *prop.Cluster) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterProp.Name,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: makeLbEndpoint(clusterProp.SocketAddress),
		}},
	}
}

func makeRoute(routeName string, clusterName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name: "local_service",
			Domains: []string{
				"*",
			},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
					},
				},
			}},
		}},
	}
}

func makeHTTPLocalRateLimit(maxTokens, tokensPerFill uint32) *local_ratelimitv3.LocalRateLimit {
	return &local_ratelimitv3.LocalRateLimit{
		StatPrefix: "http_local_rate_limiter",
		TokenBucket: &typev3.TokenBucket{
			MaxTokens: maxTokens,
			TokensPerFill: &wrappers.UInt32Value{
				Value: tokensPerFill,
			},
			FillInterval: &duration.Duration{
				Seconds: 5,
			},
		},
		FilterEnabled: &core.RuntimeFractionalPercent{
			RuntimeKey: "local_rate_limit_enabled",
			DefaultValue: &typev3.FractionalPercent{
				Numerator:   100,
				Denominator: typev3.FractionalPercent_HUNDRED,
			},
		},
		FilterEnforced: &core.RuntimeFractionalPercent{
			RuntimeKey: "local_rate_limit_enforced",
			DefaultValue: &typev3.FractionalPercent{
				Numerator:   100,
				Denominator: typev3.FractionalPercent_HUNDRED,
			},
		},
		ResponseHeadersToAdd: []*core.HeaderValueOption{
			{
				AppendAction: core.HeaderValueOption_APPEND_IF_EXISTS_OR_ADD,
				Header: &core.HeaderValue{
					Key:   "x-local-rate-limit",
					Value: "true",
				},
			},
		},
		LocalRateLimitPerDownstreamConnection: false,
	}
}

func makeHTTPListener(listenerName string, routeName string) *listener.Listener {
	routerConfig, _ := anypb.New(&router.Router{})
	localRouterConfig, _ := anypb.New(makeHTTPLocalRateLimit(1000, 1))
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: makeRoute(routeName, ClusterName),
		},

		HttpFilters: []*hcm.HttpFilter{
			{
				Name:       wellknown.HTTPRateLimit,
				ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: localRouterConfig},
			},

			{
				Name:       wellknown.Router,
				ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
			},
		},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

//func GenerateSnapshot(clusterProp prop.Cluster) *cache.Snapshot {
//	var clusterProp = prop.Cluster{
//		Name: ClusterName,
//		SocketAddress: []prop.SocketAddress{
//			{
//				Protocol:  core.SocketAddress_TCP,
//				Address:   UpstreamHost,
//				PortValue: UpstreamPort,
//			},
//			{
//				Protocol:  core.SocketAddress_TCP,
//				Address:   UpstreamHost,
//				PortValue: 8081,
//			},
//		},
//	}
//
//	snap, _ := cache.NewSnapshot(version.GetNewVersion(),
//		map[resource.Type][]types.Resource{
//			resource.ClusterType: {makeCluster(clusterProp)},
//			//resource.RouteType:    {makeRoute(RouteName, ClusterName)},
//			resource.ListenerType: {makeHTTPListener(ListenerName, RouteName)},
//		},
//	)
//	return snap
//}
