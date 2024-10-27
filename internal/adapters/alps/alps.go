package alps

import (
	"context"
	"github.com/yifeistudio-developer/canoe/internal/adapters/grpc"
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
	"github.com/yifeistudio-developer/wharf/golang/alps"
	gc "google.golang.org/grpc"
)

type Adapter struct {
}

func (a *Adapter) GetAccountPrincipals() (domain.AlpsUserProfile, error) {
	var userProfile domain.AlpsUserProfile
	err := grpc.Dialog(func(conn *gc.ClientConn) error {
		client := alps.NewAuthenticationServiceClient(conn)
		principals, err := client.GetAccountPrincipals(context.Background(), &alps.CredentialRequest{})
		if err != nil {
			return err
		}
		userProfile = domain.AlpsUserProfile{
			Username: principals.Username,
			Avatar:   principals.Avatar,
			Nickname: principals.Nickname,
		}
		return nil
	})
	if err != nil {
		return userProfile, err
	}
	return userProfile, nil
}
