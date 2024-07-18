package remote

import (
	"canoe/internal/remote/generated"
	"context"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

type Server struct {
	generated.UnimplementedAuthenticationServiceServer
}

type ServerInterceptor struct {
	Logger *golog.Logger
}

func (si *ServerInterceptor) UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	logger := si.Logger
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("Request to %s with metadata: %v", info.FullMethod, md)
	}
	resp, err = handler(ctx, req)
	if err != nil {
		logger.Info("Error serving %s: %v", info.FullMethod, err)
	} else {
		logger.Info("Response from %s", info.FullMethod)
	}

	return resp, err
}
