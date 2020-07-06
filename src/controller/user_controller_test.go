package controller_test

import (
	"context"
	"io"
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

func Test_CreateUserImage(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	userClient := proto.NewUserServiceClient(conn)
	imageClient, err := userClient.CreateUserImage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("../testdata/sample.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			log.Printf("発生ー: %v", err)
			t.Fatal(err)
		}
		imageReq := &proto.CreateUserImageRequest_ImageBytes{ImageBytes: buf}
		err = imageClient.Send(&proto.CreateUserImageRequest{Image: imageReq})
		log.Printf("っっｚ: %v\n", err)
		if err != nil {
			t.Errorf("fail to send image bytes: %v", err)
			return
		}
	}

	userID := "560dd795-bc60-4046-83c6-e9b4a06a8ef2"
	userIDReq := &proto.CreateUserImageRequest_UserId{UserId: userID}
	if err := imageClient.Send(&proto.CreateUserImageRequest{Image: userIDReq}); err != nil {
		t.Errorf("fail to send user id: %v", err)
	}

	if _, err := imageClient.CloseAndRecv(); err != nil && err != io.EOF {
		t.Fatal(err)
	}
}
