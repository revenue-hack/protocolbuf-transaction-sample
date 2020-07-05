package middleware

import (
	"context"
	"database/sql"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/revenue-hack/protobuf-transaction-sample/src/config"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

var (
	transactionMethods = []string{
		"Create",
		"Update",
		"Delete",
	}
)

func isTransactionMethod(method string) bool {
	for _, tMethod := range transactionMethods {
		if strings.HasPrefix(method, tMethod) {
			return true
		}
	}
	return false
}

func DBInterceptor(conn *sql.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fullMethods := strings.Split(info.FullMethod, "/")
		method := fullMethods[len(fullMethods)-1]
		var tx *sql.Tx
		if isTransactionMethod(method) {
			transaction, err := conn.BeginTx(ctx, nil)
			if err != nil {
				return nil, xerrors.New("transaction error")
			}
			tx = transaction

			ctx = context.WithValue(ctx, config.DBKey, transaction)
		} else {
			ctx = context.WithValue(ctx, config.DBKey, conn)
		}
		resp, err := handler(ctx, req)
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			return nil, xerrors.Errorf("fail to handle transaction: %w", err)
		}

		if tx != nil {
			if err := tx.Commit(); err != nil {
				return nil, xerrors.Errorf("fail to handle transaction commit: %w", err)
			}
		}

		return resp, nil
	}
}

func DBInterceptorForStream(conn *sql.DB) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fullMethods := strings.Split(info.FullMethod, "/")
		method := fullMethods[len(fullMethods)-1]
		var tx *sql.Tx
		ctx := ss.Context()
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		if isTransactionMethod(method) {
			transaction, err := conn.BeginTx(ss.Context(), nil)
			if err != nil {
				return xerrors.Errorf("fail to begin transaction: %w", err)
			}
			tx = transaction

			ctx = context.WithValue(ctx, config.DBKey, transaction)
		} else {
			ctx = context.WithValue(ctx, config.DBKey, conn)
		}
		wrappedStream.WrappedContext = ctx
		err := handler(srv, wrappedStream)
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			return xerrors.Errorf("fail to handle transaction: %w", err)
		}

		if tx != nil {
			if err := tx.Commit(); err != nil {
				return xerrors.Errorf("fail to handle transaction commit: %w", err)
			}
		}

		return nil
	}
}
