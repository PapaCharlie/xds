package main

import (
	"context"
	"log"
	"net"
	"time"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/xds"
)

func main() {
	go func() {
		lis, err := net.Listen("tcp", ":16464")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, new(server))
		grpcServer.Serve(lis)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := grpc.DialContext(ctx, "xds:///foo", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	<-ctx.Done()
}

type server struct {
	discovery.UnimplementedAggregatedDiscoveryServiceServer
}

func (*server) StreamAggregatedResources(stream discovery.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		if req.ResponseNonce != "" {
			// ACK, ignore
			continue
		}

		switch req.TypeUrl {
		case resource.ListenerType:
		case resource.ClusterType:
		case resource.RouteType:
		case resource.EndpointType:
		default:
			// ignore unknown type
		}

	}
}
