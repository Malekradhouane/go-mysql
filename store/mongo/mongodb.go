package mongo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

//Config db config
type Config struct {
	Host     string
	Port     int
	URI      string
	DB       string
	User     string
	Password string
	Socket   string
	Env      string
}

//func (cfg *Config) URI() string {
//	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s timezone=Europe/Paris  sslmode=disable",
//		cfg.Host,
//		cfg.Port,
//		cfg.DB,
//		cfg.User,
//		cfg.Password,
//	)
//}

////Client mongo client
//type Client struct {
//	db *mongo.Database
//}

////NewClient db client constructor
//func NewClient(c *Config, models []interface{}) (*Client, error) {
//	db, err := mongo.Connect("mongodb", c.URI())
//	if err != nil {
//		return nil, errors.Wrap(err, "gorm Open")
//	}
//
//	db = db.Debug()
//	err = db.AutoMigrate(models...).Error
//	if err != nil {
//		return nil, errors.Wrap(err, "gorm AutoMigrate")
//	}
//
//	db.DB().SetMaxIdleConns(0)
//	db.DB().SetMaxOpenConns(10)
//
//	return &Client{db: db, models: models}, nil
//}
//
//// Teardown teardown all db tables
//func (c *Client) Teardown() error {
//	return errors.Wrap(c.db.DropTableIfExists(c.models...).Error, "DropTableIfExists")
//}
//
//// Shutdown close client connection to db
//func (c *Client) Shutdown() error {
//	return errors.Wrap(c.db.Close(), "Close")
//}
//

// MongoInstance : MongoInstance Struct
type Client struct {
	Client *mongo.Client
	DB     *mongo.Database
	models []interface{}
}

// MI : An instance of MongoInstance Struct
var MI Client

// ConnectDB - database connection
func ConnectDB(models []interface{}) *Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user:dataimpactpassword@test-cdi.agxkj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected!")

	MI = Client{
		Client: client,
		DB:     client.Database(os.Getenv("DATABASE_NAME")),
		models: models,
	}
	return &MI
}

// Teardown teardown all db tables
func (c *Client) Teardown() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return errors.Wrap(c.DB.Drop(ctx), "DropTableIfExists")
}
