package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github/malekradhouane/test-cdi/encrypt"
	"github/malekradhouane/test-cdi/errs"
	"github/malekradhouane/test-cdi/store"
	"io/ioutil"
	"os"
)

//UserService user service
type UserService struct {
	userStore store.UserStore
}

//NewUserService constructs a new UserService
func NewUserService(us store.UserStore) *UserService {
	return &UserService{
		userStore: us,
	}
}

const (
	filename = "files/DataSet.json"
)

//Create a user
func (us *UserService) Create(ctx context.Context) ([]*store.User, error) {
	var users []*store.User
	// Reading file
	path, err := os.Getwd()
	if err != nil {
		return nil, err

	}
	fmt.Printf("reading %s/%s\n", path, filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err

	}

	result := []store.User{}
	err = json.Unmarshal([]byte(file), &result)
	if err != nil {
		return nil, err
	}

	for _, user := range result {
		if !us.userStore.IsEmailTaken(ctx, user.Email) {
			hashedPassword, err := encrypt.Hash(user.Password)
			if err != nil {
				return nil, err
			}
			user.Password = string(hashedPassword)
			user, err := us.userStore.CreateUser(ctx, &user)
			if err != nil {
				return nil, err
			}
			fmt.Sprintf("%s%s", user.ID, user.Email)
			filename := user.ID
			_, err = os.Stat(path + "/files/" + filename)
			if os.IsNotExist(err) {
				file, err := os.Create(path + "/files/" + filename)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				fmt.Fprintf(file, "%s ", user.Data)
			} else {
				return nil, errs.ErrFileAlreadyExist
			}

			fmt.Println("File created successfully", filename)
			users = append(users, user)
		}
	}
	return users, nil

}

//List users   _
func (us *UserService) List(ctx context.Context) ([]*store.User, error) {
	users, err := us.userStore.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

//List users   _
func (us *UserService) GetUser(ctx context.Context, id string) (*store.User, error) {
	user, err := us.userStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Delete _
func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	err := us.userStore.DeleteUser(ctx, id)
	if err == nil {
		// Reading file
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Remove(path + "/files/" + id)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

//Update a User
func (us *UserService) UpdateUser(ctx context.Context, req *store.User, id string) error {
	err := us.userStore.UpdateUser(ctx, req, id)
	if err != nil {
		return err
	}
	// Reading file
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fd, err := os.Open(path + "/files/" + id)

	if err != nil {
		return err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)

		line, _ := reader.ReadString('\n')
		if err != nil{
			return err
		}
		if len(req.Data) != 0 && req.Data != line{
			err = ioutil.WriteFile(path + "/files/" + id, []byte(req.Data), 0)
			if err != nil {
				panic(err)
			}
		}
	return nil
}