package database

import (
	"context"

	"euclia.xyz/accounts-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	NOTSCOLLECTION = "notifications"
)

// type IINotificationsDao interface {
// 	CreateNotification(not models.Notification) (notCreated models.Notification, err error)
// 	GetNotByToIDAndMail(id string, email string) (notFound models.Notification, err error)
// 	UpdateNotification(not models.Notification) (updatedNot models.Notification, err error)
// 	DeleteNotification(id string) (deletedNot int64, err error)
// }

type NotsDao struct {
}

func (notDao *NotsDao) CreateNotification(notGot models.Notification) (notCreated models.Notification, err error) {
	res, err := db.Collection(NOTSCOLLECTION).InsertOne(context.TODO(), notGot)
	if err == nil {
		id := res.InsertedID.(string)
		notGot.Id = id
	}
	return notGot, err
}

func (notDao *NotsDao) GetNotByToIDAndMail(id string, email string) (notFound models.Notification, err error) {
	filter := bson.D{primitive.E{Key: "email_to", Value: email}, primitive.E{Key: "to", Value: email}}
	var notif models.Notification
	err = db.Collection(NOTSCOLLECTION).FindOne(context.TODO(), filter).Decode(&notif)
	return notif, err
}

func (notDao *NotsDao) FindInvitation(id string, email string, invitationto string) (notFound models.Notification, err error) {
	filter := bson.D{primitive.E{Key: "email_to", Value: email}, primitive.E{Key: "to", Value: email}, primitive.E{Key: "invitation_to", Value: invitationto}}
	var notif models.Notification
	err = db.Collection(NOTSCOLLECTION).FindOne(context.TODO(), filter).Decode(&notif)
	return notif, err
}

func (notDao *NotsDao) UpdateNotification(not models.Notification) (updatedNot models.Notification, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: not.Id}}
	update := bson.M{"$set": bson.M{"viewed": not.Viewed}}
	_, erro := db.Collection(NOTSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	// var b []byte
	// resp.UnmarshalBSON(b)
	return not, erro
}

func (notDao *NotsDao) DeleteNotification(id string) (deletedNot int64, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	dr, err := db.Collection(NOTSCOLLECTION).DeleteOne(context.TODO(), filter)
	return dr.DeletedCount, err
}
