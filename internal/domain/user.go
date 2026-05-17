package domain

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("user not found")
	ErrDuplicate = errors.New("user already exists")
	ErrInvalid   = errors.New("invalid user data")
)

type User struct {
	ID   string
	Name string
	Age  int
}
