package store

import (
	"context"
	"github/malekradhouane/test-cdi/api"
)

type Store interface {
	UserStore
}

//UserStore represents the interface to manage users storage
type UserStore interface {
	CreateUser(context.Context, *User) (*User, error)
	IsEmailTaken(context.Context, string) bool
	Authenticate(context.Context, *api.Login) (*User, error)
	GetAllUsers(context.Context) ([]*User, error)
	Get(context.Context, string) (*User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context,*User, string)  error
}
