package model

import "time"

type User struct {
	tableName    struct{}   `sql:"user"`
	ID           int64      `sql:"id,pk"`
	Login        string     `sql:"login,notnull"`
	PasswordHash string     `sql:"hash,notnull"`
	Created      time.Time  `sql:"created,notnull"`
	Updated      time.Time  `sql:"updated,notnull"`
	Deleted      *time.Time `sql:"deleted"`
}

type UserFilter struct {
}
