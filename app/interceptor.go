package app

import (
	"context"
	"github.com/sanches1984/auth/app/errors"
	database "github.com/sanches1984/gopkg-pg-orm"
	dbmw "github.com/sanches1984/gopkg-pg-orm/middleware"
	"google.golang.org/grpc"
)

func databaseInterceptor(db database.IClient) grpc.UnaryServerInterceptor {
	return dbmw.NewDBServerInterceptor(
		func(ctx context.Context) database.IClient {
			db = db.WrapWithContext(ctx)
			return db
		},
	)
}

func errorConvertInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, errors.Convert(err)
		}
		return resp, nil
	}
}
