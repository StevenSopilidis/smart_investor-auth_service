package grpc_client

import (
	"context"
	"fmt"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/user_service_client/grpc_client/generated"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GrpcUserServiceClient struct {
	client generated.UserGrpcServiceClient
}

func NewGRPCUserServiceClient(addr string) (*GrpcUserServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to user_service: %s", err.Error())
	}

	client := generated.NewUserGrpcServiceClient(conn)
	return &GrpcUserServiceClient{
		client: client,
	}, nil
}

func (c *GrpcUserServiceClient) FindUserByEmail(email string) (domain.User, error) {
	res, err := c.client.FindUserByEmail(context.Background(), &generated.FindUserByEmailRequest{
		Email: email,
	})

	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return domain.User{}, &app_errors.UserNotFound{}
		} else {
			return domain.User{}, &app_errors.ServiceUnreachable{}
		}
	}

	return domain.User{
		Email:         res.GetEmail(),
		Password:      res.GetPassword(),
		EmaiLVerified: res.GetEmailVerified(),
	}, nil
}
