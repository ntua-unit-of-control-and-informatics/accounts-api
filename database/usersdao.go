package database

import (
	"context"

	"euclia.xyz/accounts-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	USERSCOLLECTION = "users"
)

type UsersDao struct{}

// Cretes User if not existing on DB
func (uDao *UsersDao) CreateUser(userGot models.User) (user models.User, err error) {
	res, err := db.Collection(USERSCOLLECTION).InsertOne(context.TODO(), userGot)
	if err == nil {
		id := res.InsertedID.(string)
		userGot.Id = id
	}
	return userGot, err
}

// Gets user by ID
func (uDao *UsersDao) GetUserById(id string) (userFound models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var user models.User
	err = db.Collection(USERSCOLLECTION).FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

// Deletes user by ID
func (uDao *UsersDao) DeleteUser(id string) (userDeleted int64, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	dr, err := db.Collection(USERSCOLLECTION).DeleteOne(context.TODO(), filter)
	return dr.DeletedCount, err
}

// Updates User
func (uDao *UsersDao) UpdateUser(user models.User) (userUpdated models.User, err error) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{
		"$set": bson.M{
			"preferred_username": user.PreferedUsername,
			"meta":               user.MetaInfo,
			"occupation":         user.Occupation,
			"occupation_at":      user.OccupationAt,
			"given_name":         user.GivenName,
			"name":               user.Name,
			"family_name":        user.FamilyName,
		},
	}
	_, erro := db.Collection(USERSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	return user, erro
}

func (uDao *UsersDao) UpdateUserNames(user models.User) (userUpdated models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.M{
		"$set": bson.M{
			"given_name":  user.GivenName,
			"name":        user.Name,
			"family_name": user.FamilyName,
			"email":       user.Email,
		},
	}
	_, erro := db.Collection(USERSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	return user, erro
}

func (uDao *UsersDao) UpdateUserName(user models.User) (userUpdated models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.M{"$set": bson.M{"preferred_username": user.PreferedUsername}}
	_, erro := db.Collection(USERSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	return user, erro
}

func (uDao *UsersDao) UpdateUserOrganizations(user models.User) (userUpdated models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.M{"$set": bson.M{"organizations": user.Organizations}}
	_, erro := db.Collection(USERSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	return user, erro
}

// Updates Users Meta
func (uDao *UsersDao) UpdateUsersSpentOn(user models.User) (userUpdated models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.M{"$set": bson.M{"meta": user.MetaInfo}}
	resp, erro := db.Collection(USERSCOLLECTION).UpdateOne(context.TODO(), filter, update)
	var b []byte
	resp.UnmarshalBSON(b)
	return user, erro
}

// Get users paginated
func (uDao *UsersDao) GetUsersPaginated(min int64, max int64) (counted int64, usersFound []*models.User, err error) {
	findOptions := options.Find()
	findOptions.SetSkip(min)
	findOptions.SetLimit(max)
	var users []*models.User
	nullOptions := options.Count()
	countedAll, err := db.Collection(USERSCOLLECTION).CountDocuments(context.TODO(), bson.D{{}}, nullOptions)
	cursor, err := db.Collection(USERSCOLLECTION).Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return counted, users, err
	}
	for cursor.Next(context.TODO()) {
		var elem models.User
		err := cursor.Decode(&elem)
		if err != nil {
			return countedAll, users, err
		}
		users = append(users, &elem)
	}
	if err := cursor.Err(); err != nil {
		cursor.Close(context.TODO())
		return countedAll, users, err
	}
	return countedAll, users, err
}

// Finds users with email
func (uDao *UsersDao) FindUserByEmail(email string) (userFound models.User, err error) {
	filter := bson.D{primitive.E{Key: "email", Value: email}}
	var user models.User
	err = db.Collection(USERSCOLLECTION).FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func (uDao *UsersDao) SearchUserByNameAndMail(name string) (counted int, orgsFound []*models.User, err error) {
	filter := bson.D{primitive.E{Key: "email", Value: ".*" + name + ".*"}, primitive.E{Key: "username", Value: ".*" + name + ".*"}}
	// update := bson.M{"$set": bson.M{"meta": user.MetaInfo}}
	cur, erro := db.Collection(USERSCOLLECTION).Find(context.TODO(), filter)
	var results []*models.User
	for cur.Next(context.TODO()) {
		var elem models.User
		cur.Decode(&elem)
		results = append(results, &elem)
	}
	found := len(results)
	// Close the cursor once finished
	cur.Close(context.TODO())
	return found, results, erro
}

func (uDao *UsersDao) FindOrgsUsers(org string) (counted int, orgsFound []*models.User, err error) {
	filter := bson.D{primitive.E{Key: "organizations", Value: org}}
	// update := bson.M{"$set": bson.M{"meta": user.MetaInfo}}
	cur, erro := db.Collection(USERSCOLLECTION).Find(context.TODO(), filter)
	var results []*models.User
	for cur.Next(context.TODO()) {
		var elem models.User
		cur.Decode(&elem)
		results = append(results, &elem)
	}
	found := len(results)
	// Close the cursor once finished
	cur.Close(context.TODO())
	return found, results, erro
}

// func CreateUser(user User) {
// 	log.Println(user)
// 	_, err := db.Collection(COLLECTION).InsertOne(context.TODO(), user)
// 	if err != nil {
// 		log.Panicf(err.Error())
// 	}
// 	// return insertResult, err
// }

// func GetUserBiId(id string, app *main.App) (user User, err error) {
// 	var userFound User

// 	var user1 = User{Email: "sd", Id: "adf", Username: "adf"}
// 	return user1, err
// 	db.Collection(COLLECTION).FindOne()(id)).One(&userFound)
// }
