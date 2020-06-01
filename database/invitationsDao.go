package database

import (
	"context"

	"euclia.xyz/accounts-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/options"
)

const (
	INVITATIONSCOLLECTION = "invitations"
)

// type IInvitationsDao interface {
// 	CreateInvitation(not models.Notification) (invCreated models.Notification, err error)
// 	GetByEmail(email string) (invFound models.Invitation, err error)
// 	GetById(id string) (invFound models.Invitation, err error)
// 	UpdateInvitation(inv models.Invitation) (invUpdated models.Invitation, err error)
// 	DeleteInvitation(inv models.Invitation) (deletedNot int64, err error)
// }

type InvitationsDao struct{}

func (invDao *InvitationsDao) CreateInvitation(invGot models.Invitation) (invCreated models.Invitation, err error) {
	res, err := db.Collection(INVITATIONSCOLLECTION).InsertOne(context.TODO(), invGot)
	if err == nil {
		id := res.InsertedID.(string)
		invGot.Id = id
	}
	return invGot, err
}

func (invDao *InvitationsDao) GetByEmail(email string, skip int64, max int64, viewed bool) (inv []models.Invitation, err error) {
	filter := bson.D{primitive.E{Key: "email_to", Value: email}, primitive.E{Key: "viewed", Value: viewed}}
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(max)
	var invit []models.Invitation
	cursor, err := db.Collection(INVITATIONSCOLLECTION).Find(context.TODO(), filter, findOptions)
	cursor.All(context.TODO(), &invit)
	return invit, err
}

func (invDao *InvitationsDao) CountWithEmail(email string, viewed bool) (total int64, err error) {
	filter := bson.D{primitive.E{Key: "email_to", Value: email}, primitive.E{Key: "viewed", Value: viewed}}
	totalF, errIf := db.Collection(INVITATIONSCOLLECTION).CountDocuments(context.TODO(), filter)
	return totalF, errIf
}

func (invDao *InvitationsDao) GetReceivedById(id string) (invFound models.Invitation, err error) {
	filter := bson.D{primitive.E{Key: "to", Value: id}}
	var invit models.Invitation
	err = db.Collection(INVITATIONSCOLLECTION).FindOne(context.TODO(), filter).Decode(&invit)
	return invit, err
}

func (invDao *InvitationsDao) GetSentById(id string) (invFound models.Invitation, err error) {
	filter := bson.D{primitive.E{Key: "from", Value: id}}
	var invit models.Invitation
	err = db.Collection(INVITATIONSCOLLECTION).FindOne(context.TODO(), filter).Decode(&invit)
	return invit, err
}

func (invDao *InvitationsDao) UpdateInvitation(inv models.Invitation) (invUpdated models.Invitation, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: inv.Id}}
	update := bson.M{"$set": bson.M{"body": inv.Body, "viewed": inv.Viewed, "action": inv.Action}}
	_, erro := db.Collection(INVITATIONSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	return inv, erro
}

func (invDao *InvitationsDao) DeleteInvitation(id string) (deletedNot int64, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	dr, err := db.Collection(INVITATIONSCOLLECTION).DeleteOne(context.TODO(), filter)
	return dr.DeletedCount, err
}
