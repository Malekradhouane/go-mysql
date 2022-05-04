// Package classification of todolist API
//
// Documentation for todolist API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package route

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github/malekradhouane/test-cdi/errs"
	"github/malekradhouane/test-cdi/service"
	. "github/malekradhouane/test-cdi/store"
	"net/http"
)

//TodoActions represents users controller actions
type TodoActions struct {
	todoService *service.TodoService
}

//NewTodoActions constructor
func NewTodoActions(us *service.TodoService) *TodoActions {
	return &TodoActions{
		todoService: us,
	}
}

// swagger:route POST /todolist todolist create
// Create a new todolist
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse
//CreateUser creates a new todo
func (as TodoActions) Create(c *gin.Context) {
	req := new(Todo)

	err := c.ShouldBindBodyWith(req, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_, err = as.todoService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo added"})

}

// swagger:route GET /todolists todoList
// Return a list of todos from the database
// responses:
//	200: todoList
func (as TodoActions) Get(c *gin.Context) {
	todoList, err := as.todoService.GetAll(c.Request.Context())
	if err == errs.ErrNoSuchEntity {
		c.JSON(http.StatusOK, gin.H{"message": "Todolist not found"})
		return
	}

	c.JSON(http.StatusOK, todoList)
}

// swagger:route DELETE /todolists/{id} todolist Delete
// delete a todolist
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (ua TodoActions) Delete(c *gin.Context) {
	id := c.Param("id")
	err := ua.todoService.Delete(c.Request.Context(), id)
	if err != nil {
		if err == errs.ErrNoSuchEntity {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No user with id : %s", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted with success"})
}
