package mongo

import (
	"context"
	"github.com/chidiwilliams/flatbson"
	"github.com/jinzhu/gorm"
	"github/malekradhouane/test-cdi/api"
	"github/malekradhouane/test-cdi/encrypt"
	"github/malekradhouane/test-cdi/errs"
	. "github/malekradhouane/test-cdi/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//CreateUser create a user
func (c *Client) CreateUser(ctx context.Context, user *User) (*User, error) {
	_, err := c.Client.Database("user").Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Authenticate retrieve a user
func (c *Client) Authenticate(ctx context.Context, login *api.Login) (*User, error) {
	user := new(User)
	filter := bson.D{{"email", login.Email}}
	err := c.Client.Database("user").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errs.ErrNoSuchEntity
		}
		return nil, err
	}
	err = encrypt.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//IsEmailTaken retrieve a user
func (c *Client) IsEmailTaken(ctx context.Context, email string) bool {
	filter := bson.D{{"email", email}}
	cursor, err := c.Client.Database("user").Collection("users").CountDocuments(ctx, filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	return cursor != 0
}


func (c *Client) GetAllUsers(ctx context.Context) ([]*User, error) {
	var Users []*User  //slice for multiple documents
	cursor, err := c.Client.Database("user").Collection("users").Find(ctx, bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var user *User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}


		Users = append(Users, user) // appending document pointed by Next()
}
	cursor.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	return Users, err
}

//Get retrieve a user
func (c *Client) Get(ctx context.Context, id string) (*User, error) {
	user := new(User)
	filter := bson.D{{"id", id}}
	err := c.Client.Database("user").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errs.ErrNoSuchEntity
		}
		return nil, err
	}
	return user, nil
}

//DeleteUser deletes a user with given ID
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	filter := bson.D{{"id", id}}
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	_, err :=  c.Client.Database("user").Collection("users").DeleteOne(ctx, filter, opts)
	if err != nil {
		return err
	}
	return nil
}

//Update User
func (c *Client) UpdateUser(ctx context.Context, user *User, id string)  error {
	filter := bson.D{{"id", id}}
	u := &User{}
	us, err := flatbson.Flatten(user)
	update := bson.M{
		"$set": us,
	}
	 err = c.Client.Database("user").Collection("users").FindOneAndUpdate(ctx, filter, update).Decode(u)
	if err != nil {
		if errs.IsNoSuchEntityError(err) {
			return errs.ErrNoSuchEntity
		}
		return err
	}
	return nil
}

