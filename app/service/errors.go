package service

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserIsDeleted = errors.New("user is deleted")
var ErrIncorrectPassword = errors.New("incorrect password")
