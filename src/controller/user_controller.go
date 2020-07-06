package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

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

func (s *CreateUserController) CreateUserImage(stream proto.UserService_CreateUserImageServer) error {
	var buf bytes.Buffer
	var userID string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, xerrors.Errorf("fail to get stream request: %w", err).Error())
		}

		switch input := in.Image.(type) {
		case *proto.CreateUserImageRequest_ImageBytes:
			if _, err := buf.Write(input.ImageBytes); err != nil {
				return xerrors.Errorf("fail to read image bytes: %w", err)
			}
		case *proto.CreateUserImageRequest_UserId:
			userID = input.UserId
		default:
			break
		}
	}

	file, err := os.Create(fmt.Sprintf("./%s.jpg", userID))
	if err != nil {
		return xerrors.Errorf("fail to create: %w", err)
	}
	if _, err := io.Copy(file, &buf); err != nil {
		return xerrors.Errorf("fail to copy: %w", err)
	}

	return nil
}
