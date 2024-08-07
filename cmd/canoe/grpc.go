package main

import (
	"canoe/internal/remote"
	"canoe/internal/remote/generated"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	"net"
)

func startGrpcServer(log *golog.Logger) *grpc.Server {
	interceptor := remote.ServerInterceptor{Logger: log}
	gRpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor))
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		generated.RegisterAuthenticationServiceServer(gRpcServer, &remote.Server{})
		for k, v := range gRpcServer.GetServiceInfo() {
			log.Info("service info: ", k, v)
		}
		if err := gRpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return gRpcServer
}
