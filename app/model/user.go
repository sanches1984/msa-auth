package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	tableName    struct{}   `sql:"users"`
	ID           int64      `sql:"id,pk"`
	Login        string     `sql:"login,notnull"`
	PasswordHash string     `sql:"hash,notnull"`
	Created      time.Time  `sql:"created,notnull"`
	Updated      time.Time  `sql:"updated,notnull"`
	Deleted      *time.Time `sql:"deleted"`
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
