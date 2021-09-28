package route

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github/malekradhouane/test-cdi/auth"
	"github/malekradhouane/test-cdi/errs"
	"github/malekradhouane/test-cdi/service"
	"github/malekradhouane/test-cdi/store"
	"net/http"
)

//UserActions represents users controller actions
type UserActions struct {
	userService *service.UserService
}

//NewUserActions constructor
func NewUserActions(us *service.UserService) *UserActions {
	return &UserActions{
		userService: us,
	}
}

//CreateUser creates a new user
func (as UserActions) CreateUser(c *gin.Context) {
	users, err := as.userService.Create(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	if len(users) != 0 {
		c.JSON(http.StatusCreated, users)
	} else {
		c.JSON(http.StatusAccepted, "no user to add")
	}
}

func (as UserActions) ListUsers(c *gin.Context) {
	authed := auth.ExtractAuthenticated(c)
	if authed == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Identification info required for getting users"})
		return
	}
	requests, err := as.userService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

//GetUser _
func (as UserActions) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := as.userService.GetUser(c.Request.Context(), id)
	if err == errs.ErrNoSuchEntity {
		c.JSON(http.StatusOK, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ua UserActions) DeleteUser(c *gin.Context) {
	authed := auth.ExtractAuthenticated(c)
	if authed == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Identification info required for deleting user"})
		return
	}
	id := c.Param("id")
	err := ua.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if err == errs.ErrNoSuchEntity {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No user with id : %s", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted with success"})
}


//UpdateUser
func (ua UserActions) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	req := new(store.User)
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
	err = ua.userService.UpdateUser(c.Request.Context(), req, id)
	if err != nil {
		if errs.IsNoSuchEntityError(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "user updated"})

}
