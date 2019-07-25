package database

import (
	"context"
	"quots/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION = "users"
)

type UsersDao struct{}

// Cretes User if not existing on DB
func (uDao *UsersDao) CreateUser(userGot models.User) (user models.User, err error) {
	res, err := db.Collection(COLLECTION).InsertOne(context.TODO(), userGot)
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
	err = db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

// Updates Users Credits
func (uDao *UsersDao) UpdateUserCredits(user models.User) (userUpdated models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.M{"$set": bson.M{"credits": user.Credits}}
	resp, erro := db.Collection(COLLECTION).UpdateOne(context.TODO(), filter, update)
	var b []byte
	resp.UnmarshalBSON(b)
	return user, erro
}

// Updates Users Spent
func (uDao *UsersDao) UpdateUsersSpentOn(user models.User) (userUpdated models.User, err error) {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Id}}
	update := bson.M{"$set": bson.M{"spenton": user.Spenton}}
	resp, erro := db.Collection(COLLECTION).UpdateOne(context.TODO(), filter, update)
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
	countedAll, err := db.Collection(COLLECTION).CountDocuments(context.TODO(), bson.D{{}}, nullOptions)
	cursor, err := db.Collection(COLLECTION).Find(context.TODO(), bson.D{{}}, findOptions)
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
	err = db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	return user, err
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
