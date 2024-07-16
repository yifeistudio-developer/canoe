package remote

import (
	"canoe/internal/remote/generated"
	"context"
)

type Server struct {
	generated.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *generated.HelloRequest) (*generated.HelloReply, error) {
	return &generated.HelloReply{Message: "hello " + in.Name}, nil
}
