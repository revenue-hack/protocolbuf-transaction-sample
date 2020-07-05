package main

import (
	"context"

	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/rdb"
	"github.com/revenue-hack/protobuf-transaction-sample/src/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateUserController struct{}

func NewCreateUserController() *CreateUserController {
	return &CreateUserController{}
}

func (s *CreateUserController) CreateUser(ctx context.Context, in *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	conn, err := rdb.ExecFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "fail to get RDB Connection")
	}

	userRepo := repoimpl.NewUserRepoImpl(conn)

	app := userapp.NewCreateUserImageApp(stImpl, userRepo)
	// @todo 後でextensionも変更する
	if err := app.Exec(ctx, userID, &buf, "jpg"); err != nil {
		return err
	}

	return &proto.CreateUserResponse{}, nil
}
