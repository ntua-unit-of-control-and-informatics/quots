package database

import (
	"context"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	models "quots/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var defaultCredits int64

type E struct {
	Key   string
	Value interface{}
}

// type DB struct {
// 	db *mongo.Database
// }

type Datastore interface {
	FindById(object *interface{}) (interface{}, error)
	FindAll(object *interface{}) (interface{}, error)
	UpdateOne(object *interface{}) (interface{}, error)
	CreateOne(object *interface{}) (interface{}, error)
	DeleteOne(object *interface{}) (interface{}, error)
}

func NewDB() {
	log.Println("Starting DB")
	mongoURL := os.Getenv("MONGO_URL")
	defaultCreditsGot := os.Getenv("CREDITS")
	if defaultCreditsGot == "" {
		defaultCredits = 20
	} else {
		defC, erro := strconv.ParseInt(defaultCreditsGot, 10, 64)
		if erro != nil {
			log.Fatal(erro.Error())
		}
		defaultCredits = defC
	}
	database := os.Getenv("DATABASE")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}
	log.Println("Starting at " + mongoURL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if database == "" {
		database = "quots"
	}
	db = client.Database(database)
	if err != nil {
		log.Panicln(err.Error())
	}
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func SaveJWTSignature(signature string) (res *mongo.InsertOneResult, err error) {
	secret := models.JWTSecret{
		Signature: signature,
	}
	resp, erro := db.Collection("secret").InsertOne(context.TODO(), secret)
	return resp, erro
}

func GetJWTSignature() (secrets []*models.JWTSecret, err error) {
	var secret []*models.JWTSecret
	findOptions := options.Find()
	cursor, erro := db.Collection("secret").Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return secret, erro
	}
	for cursor.Next(context.TODO()) {
		var elem models.JWTSecret
		err := cursor.Decode(&elem)
		if err != nil {
			return secret, erro
		}
		secret = append(secret, &elem)
	}
	if err := cursor.Err(); err != nil {
		cursor.Close(context.TODO())
		return secret, erro
	}
	return secret, erro
}

// func CreateOne(entity *interface{}) (entityCreated *mongo.InsertOneResult, err error) {
// 	collection := getType(entity)
// 	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), entity)
// 	if err != nil {
// 		log.Panicf(err.Error())
// 	}
// 	return insertResult, err
// }

// func GetUsers(start int32, max int32) (entityCreated *mongo.Cursor, err error) {
// 	findOptions := options.Find()
// 	// findOptions.SetLimit(2)
// 	findOptions.SetMin(start)
// 	findOptions.SetMax(max)
// 	insertResult, err := DB.db.Collection(COLLECTION).Find(context.TODO(), bson.D{{}}, findOptions)
// 	if err != nil {
// 		log.Panicf(err.Error())
// 	}
// 	return insertResult, err
// }
