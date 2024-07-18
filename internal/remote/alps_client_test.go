package remote

import (
	"canoe/internal/remote/generated"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

func TestGetAccountPrincipals(t *testing.T) {
	conn, err := grpc.Dial("192.168.18.207:49395", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}(conn)
	client := generated.NewAuthenticationServiceClient(conn)
	r, err := client.GetAccountPrincipals(context.Background(), &generated.TicketRequest{Ticket: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	println(r)
}
