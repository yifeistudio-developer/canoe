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

func GetUserProfile(accessToken string) (*model.AlpsUserProfile, error) {
	// 获取从nacos获取实例
	instance, err := config.GetService("alps")
	if err != nil {
		return nil, err
	}
	metadata := instance.Metadata
	ip := instance.Ip
	port := metadata["gRPC_port"]
	target := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			// close error
		}
	}(conn)
	client := generated.NewAuthenticationServiceClient(conn)
	result, err := client.GetAccountPrincipals(context.Background(), &generated.CredentialRequest{Credential: accessToken,
		Type: generated.CredentialType_ACCESS_TOKEN})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	principals := result.Principals
	return &model.AlpsUserProfile{
		Username: principals.Username,
		Nickname: principals.Nickname,
		Avatar:   principals.Avatar,
	}, err
}
