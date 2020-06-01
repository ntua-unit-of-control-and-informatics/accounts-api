package database

import (
	"context"

	"euclia.xyz/accounts-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ORGSCOLLECTION = "organizations"
)

type OrgsDao struct{}

// Creates Organization
func (oDao *OrgsDao) CreateOrganization(orgGot models.Organization) (orgCreated models.Organization, err error) {
	res, err := db.Collection(ORGSCOLLECTION).InsertOne(context.TODO(), orgGot)
	if err == nil {
		id := res.InsertedID.(string)
		orgGot.Id = id
	}
	return orgGot, err
}

// Gets Organization by ID
func (uDao *OrgsDao) GetOrgById(id string) (orgGot models.Organization, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var org models.Organization
	err = db.Collection(ORGSCOLLECTION).FindOne(context.TODO(), filter).Decode(&org)
	return org, err
}

// Gets Organization by admin
func (uDao *OrgsDao) GetOrgByAdmin(id string) (counted int, orgsFound []*models.Organization, err error) {
	filter := bson.D{primitive.E{Key: "admins", Value: id}}
	cur, erro := db.Collection(ORGSCOLLECTION).Find(context.TODO(), filter)
	var results []*models.Organization
	for cur.Next(context.TODO()) {
		var elem models.Organization
		cur.Decode(&elem)
		results = append(results, &elem)
	}
	found := len(results)
	// Close the cursor once finished
	cur.Close(context.TODO())
	return found, results, erro
}

// Deletes Organization by ID
func (uDao *OrgsDao) DeleteOrg(id string) (orgsDeleted int64, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	dr, err := db.Collection(ORGSCOLLECTION).DeleteOne(context.TODO(), filter)
	return dr.DeletedCount, err
}

// Updates Organization
func (uDao *OrgsDao) UpdateOrg(orgToUpdate models.Organization) (orgUpdated models.Organization, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: orgToUpdate.Id}}
	update := bson.M{"$set": bson.M{"admins": orgToUpdate.Admins, "meta": orgToUpdate.MetaInfo, "users": orgToUpdate.Users, "markdown": orgToUpdate.Markdown, "credits": orgToUpdate.Credits, "contact": orgToUpdate.Contact, "contact_types": orgToUpdate.ContactTypes}}
	// resp, erro := db.Collection(COLLECTION).UpdateOne(context.TODO(), filter, update)
	resp, erro := db.Collection(ORGSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	var b []byte
	resp.UnmarshalBSON(b)
	return orgToUpdate, erro
}

func (uDao *OrgsDao) UpdateOrgAdmins(orgToUpdate models.Organization) (orgUpdated models.Organization, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: orgToUpdate.Id}}
	update := bson.M{"$set": bson.M{"admins": orgToUpdate.Admins}}
	// resp, erro := db.Collection(COLLECTION).UpdateOne(context.TODO(), filter, update)
	resp, erro := db.Collection(ORGSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	var b []byte
	resp.UnmarshalBSON(b)
	return orgToUpdate, erro
}

// Finds partial matching orgs
func (uDao *OrgsDao) SearchOrgByName(name string) (counted int, orgsFound []*models.Organization, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ".*" + name + ".*"}}
	// update := bson.M{"$set": bson.M{"meta": user.MetaInfo}}
	cur, erro := db.Collection(ORGSCOLLECTION).Find(context.TODO(), filter)
	var results []*models.Organization
	for cur.Next(context.TODO()) {
		var elem models.Organization
		cur.Decode(&elem)
		results = append(results, &elem)
	}
	found := len(results)
	// Close the cursor once finished
	cur.Close(context.TODO())
	return found, results, erro
}
