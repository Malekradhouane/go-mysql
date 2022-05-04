package service

import (
	"context"
	"github/malekradhouane/test-cdi/store"
	"strconv"
)

//TodoService Todoservice service
type TodoService struct {
	TodoStore store.TodoStore
}

//NewTodoService constructs a new TodoService
func NewTodoService(us store.TodoStore) *TodoService {
	return &TodoService{
		TodoStore: us,
	}
}



func (us *TodoService) Create(ctx context.Context, req *store.Todo) (*store.Todo, error) {
	Todo, err := us.TodoStore.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return Todo, nil

}



//List TodoList   _
func (us *TodoService) GetAll(ctx context.Context) ([]*store.Todo, error) {
	Todo, err := us.TodoStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return Todo, nil
}

//Delete _
func (us *TodoService) Delete(ctx context.Context, id string) error {
	requestId, _ := strconv.Atoi(id)
	err:= us.TodoStore.Delete(ctx, requestId)
	if err!= nil {
		return err
	}
	return nil
}
