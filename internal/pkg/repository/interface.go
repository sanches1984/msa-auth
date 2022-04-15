package repository

import (
	"context"
	"github.com/sanches1984/gopkg-pg-orm/repository/dao"
	"github.com/sanches1984/gopkg-pg-orm/repository/opt"
)

type ORM interface {
	FindOne(ctx context.Context, receiver interface{}, opts []opt.FnOpt) error
	FindList(ctx context.Context, receiver interface{}, opts []opt.FnOpt) error
	Insert(ctx context.Context, rec ...interface{}) error
	Update(ctx context.Context, rec interface{}, columns ...string) error
	SoftDelete(ctx context.Context, rec dao.DeletedSetter) error
	HardDeleteWhere(ctx context.Context, rec interface{}, opts []opt.FnOpt) error
}
