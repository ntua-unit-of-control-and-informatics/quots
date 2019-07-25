package database

import (
	"context"
	"quots/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	APPCOLLECTION = "applications"
)

type ApplicationDao struct{}

// Cretes User if not existing on DB
func (applicationDao *ApplicationDao) CreateApplication(applicationGot models.Application) (application models.Application, err error) {
	res, err := db.Collection(APPCOLLECTION).InsertOne(context.TODO(), applicationGot)
	if err == nil {
		id := res.InsertedID.(string)
		applicationGot.Id = id
	}
	return applicationGot, err
}

// Gets application by ID
func (applicationDao *ApplicationDao) GetApplicationBiId(id string) (appFound models.Application, err error) {
	// filter := bson.D{{"_id", id}}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var application models.Application
	err = db.Collection(APPCOLLECTION).FindOne(context.TODO(), filter).Decode(&application)
	return application, err
}

// Gets application by ID
func (applicationDao *ApplicationDao) DeleteApplicationBiId(id string) (deleted int64, err error) {
	// filter := bson.D{{"_id", id}}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	deleteResult, err := db.Collection(APPCOLLECTION).DeleteOne(context.TODO(), filter)
	return deleteResult.DeletedCount, err
}

// Updates Application
func (applicationDao *ApplicationDao) UpdateApp(applicationToUpdate models.Application) (application models.Application, err error) {
	id := applicationToUpdate.Id
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$set": bson.M{"usagescost": applicationToUpdate.UsagesCost, "usagetypes": applicationToUpdate.UsageTypes, "clientsecret": applicationToUpdate.AppSecret, "baseURLS": applicationToUpdate.BaseURLS}}
	resp, erro := db.Collection(APPCOLLECTION).UpdateOne(context.TODO(), filter, update)
	var b []byte
	resp.UnmarshalBSON(b)
	return applicationToUpdate, erro
}

// func (applicationDao *ApplicationDao) GetAllApps(min int64, max int64) (applications []*models.Application, err error) {
// 	findOptions := options.Find()
// 	findOptions.SetMin(min)
// 	findOptions.SetLimit(max)
// 	var apps []*models.Application
// 	cursor, err := db.Collection(APPCOLLECTION).Find(context.TODO(), bson.D{{}}, findOptions)
// 	defer cursor.Close(context.TODO())
// 	if err != nil {
// 		return apps, err
// 	}
// 	for cursor.Next(context.TODO()) {
// 		// create a value into which the single document can be decoded
// 		var elem models.Application
// 		err := cursor.Decode(&elem)
// 		if err != nil {
// 			return apps, err
// 		}
// 		apps = append(apps, &elem)
// 	}
// 	if err := cursor.Err(); err != nil {
// 		return apps, err
// 	}
// 	return apps, err
// }

func (applicationDao *ApplicationDao) GetAllApps(min int64, max int64) (counted int64, applications []*models.Application, err error) {
	findOptions := options.Find()
	findOptions.SetSkip(min)
	findOptions.SetLimit(max)
	// findOptions.SetLimit(2)
	var apps []*models.Application
	nullOptions := options.Count()
	countedAll, err := db.Collection(APPCOLLECTION).CountDocuments(context.TODO(), bson.D{{}}, nullOptions)
	cursor, err := db.Collection(APPCOLLECTION).Find(context.TODO(), bson.D{{}}, findOptions)
	// defer cursor.Close(context.TODO())
	if err != nil {
		return counted, apps, err
	}
	for cursor.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.Application
		err := cursor.Decode(&elem)
		if err != nil {
			return countedAll, apps, err
		}
		apps = append(apps, &elem)
	}
	if err := cursor.Err(); err != nil {
		cursor.Close(context.TODO())
		return countedAll, apps, err
	}
	return countedAll, apps, err

}
