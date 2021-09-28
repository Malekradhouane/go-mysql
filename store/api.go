package store

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

//User represents User
type User struct {
	ID         string   `bson:"id,omitempty"`
	Email      string   `bson:"email,omitempty"`
	Password   string   `bson:"password,omitempty"`
	IsActive   bool     `bson:"isActive,omitempty" default:"false"`
	Balance    string   `bson:"balance,omitempty"`
	Age        int      `bson:"age,omitempty"`
	Name       string   `bson:"name,omitempty"`
	Gender     string   `bson:"gender,omitempty"`
	Company    string   `bson:"company,omitempty"`
	Phone      string   `bson:"phone,omitempty"`
	Address    string   `bson:"address,omitempty"`
	About      string   `bson:"about,omitempty"`
	Registered string   `bson:"registered,omitempty"`
	Latitude   float32  `bson:"latitude,omitempty"`
	Longitude  float32  `bson:"longitude,omitempty"`
	Friends    []Friend `bson:"friends,omitempty"`
	Tags       []string `bson:"tags,omitempty"`
	Data       string   `bson:"data,omitempty"`
}

type Friend struct {
	Id   int    `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}

//Value marshalls to JSON value
func (r User) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *User) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, r)
	}

	return errors.New(fmt.Sprint("Failed to unmarshal JSON from the database", src))
}