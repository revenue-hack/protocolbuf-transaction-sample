package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/rdb"
	"github.com/revenue-hack/protobuf-transaction-sample/src/proto"
	"golang.org/x/xerrors"
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

	ins, err := conn.Prepare("insert into users (id, name) values (?, ?)")
	if err != nil {
		return nil, status.Error(codes.Internal, xerrors.Errorf("fail to prepare sql: %w", err).Error())
	}
	id := uuid.New().String()
	if _, err := ins.Exec(id, in.Name); err != nil {
		return nil, status.Error(codes.Internal, xerrors.Errorf("fail to execute: %w", err).Error())
	}

	return &proto.CreateUserResponse{Id: id, Name: in.Name}, nil
}
