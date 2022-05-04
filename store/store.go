package store

import (
	"context"
)


//TodoStore represents the interface to manage TodoList storage
type TodoStore interface {
	Create(context.Context, *Todo) (*Todo, error)
	GetAll(context.Context) ([]*Todo, error)
	Delete(context.Context, int) error
}
