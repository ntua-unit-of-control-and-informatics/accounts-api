package database

import (
	"context"
	"log"
	"os"

	"reflect"
	"time"

	// models "euclia.xyz/accounts-api/models"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

type E struct {
	Key   string
	Value interface{}
}

type IDatastore interface {
	FindById(T reflect.Type, id string) (interface{}, error)
	// FindAll(object *interface{}) (interface{}, error)
	// UpdateOne(object *interface{}) (interface{}, error)
	CreateOne(object *interface{}) (interface{}, error)
	// DeleteOne(object *interface{}) (interface{}, error)
}

func NewDB() {
	log.Println("Starting DB")
	mongoURL := os.Getenv("MONGO_URL")
	database := os.Getenv("DATABASE")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}
	log.Println("Starting at " + mongoURL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if database == "" {
		database = "accounts-api"
	}
	db = client.Database(database)
	if err != nil {
		log.Panicln(err.Error())
	}
}

type Datastore struct{}

func (datastore *Datastore) CreateOne(object *interface{}) (interface{}, error) {
	// objectToUpdate := reflect.ValueOf(&object).Elem()
	// objectToUpdate.CanSet()
	// t := objectToUpdate.Elem().Type().Elem()
	_, err := db.Collection(reflect.TypeOf(object).String()).InsertOne(context.TODO(), object)
	if err == nil {
		log.Println(err.Error())
		// id := inserted.InsertedID.(string)
		// userGot.Id = id
	}
	return object, err
}

func (datastore *Datastore) FindById(T reflect.Type, id string) (interface{}, error) {
	r := reflect.ValueOf(&T).Elem()
	// objectToUpdate.CanSet()
	t := r.Elem().Type().Elem()
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := db.Collection(reflect.TypeOf(T).String()).FindOne(context.TODO(), filter).Decode(&t)
	if err == nil {
		// log.Println(err.Error())
		// id := inserted.InsertedID.(string)
		// userGot.Id = id
	}
	return t, err
}
