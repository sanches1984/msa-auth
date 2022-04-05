package model

import (
	"context"
	"github.com/sanches1984/gopkg-pg-orm/repository/opt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserList []*User
type UserOrder int

const (
	UserOrderCreatedAsc UserOrder = iota
	UserOrderCreatedDesc
	UserOrderLoginAsc
	UserOrderLoginDesc
)

type User struct {
	tableName    struct{}   `pg:"users"`
	ID           int64      `pg:"id,pk"`
	Login        string     `pg:"login,notnull"`
	PasswordHash string     `pg:"password_hash,notnull"`
	Created      time.Time  `pg:"created,notnull"`
	Updated      time.Time  `pg:"updated,notnull"`
	Deleted      *time.Time `pg:"deleted"`
}

type UserFilter struct {
	ID          int64
	Login       string
	Order       UserOrder
	ShowDeleted bool
}

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.Created = time.Now()
	u.Updated = time.Now()
	return ctx, nil
}

func (u *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.Updated = time.Now()
	return ctx, nil
}

func (u *User) SetDeleted(t time.Time) {
	now := time.Now()
	u.Deleted = &now
}

func (u *User) SetHashByPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) IsPasswordCorrect(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (uo UserOrder) GetOptFn() opt.FnOpt {
	switch uo {
	case UserOrderCreatedDesc:
		return opt.Desc("id")
	case UserOrderLoginAsc:
		return opt.Asc("login")
	case UserOrderLoginDesc:
		return opt.Desc("login")
	default:
		return opt.Asc("id")
	}
}
