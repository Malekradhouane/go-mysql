package mysql

import (
	"context"
	. "github/malekradhouane/test-cdi/store"
)

//CreateUser create a user
func (c *Client) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	_, err := c.db.Exec("INSERT INTO todo (id, description, completed) VALUES (:id, :desciption, :completed)", todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (c *Client) GetAll(ctx context.Context) ([]*Todo, error) {
	var todolist []*Todo //slice for multiple documents

	return todolist, nil
}

//Delete deletes a todo  with given ID
func (c *Client) Delete(ctx context.Context, id int) error {
	_, err := c.db.Prepare("delete from todo where id=" + string(id))
	if err != nil {
		return err
	}
	return nil
}
