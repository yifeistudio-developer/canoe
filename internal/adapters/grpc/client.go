package grpc

import (
	"context"
	"github.com/yifeistudio-developer/canoe/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Printf("gRPC Error - Code: %v, Message: %v\n", st.Code(), st.Message())
		} else {
			log.Printf("Non-gRPC Error: %v\n", err)
		}
	}
	return err
}

func Dialog(executeFunc func(conn *grpc.ClientConn) error) error {
	var options []grpc.DialOption
	options = append(options,
		grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(config.GetGrpcServiceUrl(), options...)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("gRPC Error - %v\n", err)
			return
		}
	}(conn)
	if err != nil {
		log.Printf("gRPC Error - %v\n", err)
	}
	return executeFunc(conn)
}
