package controller_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/revenue-hack/protobuf-transaction-sample/src/controller"
	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/middleware"
	"github.com/revenue-hack/protobuf-transaction-sample/src/infra/rdb"
	"github.com/revenue-hack/protobuf-transaction-sample/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

const bufSize = 1024 * 1024

func setup() {
	lis = bufconn.Listen(bufSize)
	conn, err := rdb.NewConnection()
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(
		grpcmiddleware.WithStreamServerChain(
			middleware.DBInterceptorForStream(conn),
		),
		grpcmiddleware.WithUnaryServerChain(
			middleware.DBInterceptor(conn),
		),
	)
	userControl := controller.NewCreateUserController()
	proto.RegisterUserServiceServer(s, userControl)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func Test_CreateUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	userClient := proto.NewUserServiceClient(conn)
	expectedName := "testname"
	resp, err := userClient.CreateUser(ctx, &proto.CreateUserRequest{Name: expectedName})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != expectedName {
		t.Errorf("name is invalid. result is %s, expected is %s", resp.Name, expectedName)
	}
}
