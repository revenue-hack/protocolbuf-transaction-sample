package main

import (
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/middleware"
	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/rdb"
	"github.com/revenue-hack/protobuf-transaction-sample/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func recoveryFunc(p interface{}) error {
	return grpc.Errorf(codes.Internal, "Unexpected error")
}

func initService(s *grpc.Server) {
	userControl := NewCreateUserController()
	proto.RegisterUserServiceServer(s, userControl)
}

func main() {
	opts := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(recoveryFunc),
	}
	conn, err := rdb.NewConnection()
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpcmiddleware.WithStreamServerChain(
			middleware.DBInterceptorForStream(conn),
		),
		grpcmiddleware.WithUnaryServerChain(
			grpcrecovery.UnaryServerInterceptor(opts...),
			middleware.DBInterceptor(conn),
		),
	)

	initService(server)

	listenPort, err := net.Listen("tcp", ":6565")
	if err != nil {
		panic(err)
	}
	if err := server.Serve(listenPort); err != nil {
		panic(err)
	}
}
