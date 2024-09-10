package remote

import (
	"canoe/internal/config"
	"canoe/internal/model"
	"canoe/internal/remote/generated"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetUserProfile(accessToken string) *model.AlpsUserProfile {
	// 获取从nacos获取实例
	service, err := config.GetService("alps")
	if err != nil {
		fmt.Println(service, err)
	}
	instance := service.Hosts[0]
	metadata := instance.Metadata
	ip := instance.Ip
	port := metadata["gRPC_port"]
	target := fmt.Sprintf("%s:%s", ip, port)
	fmt.Println(target)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	client := generated.NewAuthenticationServiceClient(conn)
	result, err := client.GetAccountPrincipals(context.Background(), &generated.TicketRequest{Ticket: accessToken})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	principals := result.Principals
	return &model.AlpsUserProfile{
		Username: principals.Username,
		Nickname: principals.Nickname,
		Avatar:   principals.Avatar,
	}
}
