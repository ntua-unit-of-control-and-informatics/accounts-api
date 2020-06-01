package database

import (
	"context"

	"euclia.xyz/accounts-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	REQUESTSCOLLECTION = "requests"
)

// GetRequest(userId string) (reqCreated models.Request, err error)

type RequestsDao struct{}

func (reqDao *RequestsDao) CreateRequest(reqGot models.Request) (reqCreated models.Request, err error) {
	res, err := db.Collection(REQUESTSCOLLECTION).InsertOne(context.TODO(), reqGot)
	if err == nil {
		id := res.InsertedID.(string)
		reqGot.Id = id
	}
	return reqGot, err
}

func (reqDao *RequestsDao) GetRequest(userId string) (reqCreated models.Request, err error) {
	filter := bson.D{primitive.E{Key: "from_user", Value: userId}}
	var req models.Request
	err = db.Collection(REQUESTSCOLLECTION).FindOne(context.TODO(), filter).Decode(&req)
	return req, err
}
