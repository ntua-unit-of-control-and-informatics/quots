package config

import (
	"context"
	"log"
	"os"
	"time"

	oidc "github.com/coreos/go-oidc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigInterface interface {
}

type Config struct {
	MongoURL      string
	MongoClient   *mongo.Client
	Database      *mongo.Database
	TokenVerifier *oidc.IDTokenVerifier
}

func Init() (conf *Config) {
	config := new(Config)
	mongoURL := os.Getenv("MONGO_URL")
	database := os.Getenv("DATABASE")
	oidc_provider := os.Getenv("OIDC_PROVIDER")
	oidcClientId := os.Getenv("OIDC_CLIENT")
	if mongoURL == "" {
		config.MongoURL = "mongodb://localhost:27017"
	} else {
		config.MongoURL = mongoURL
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		config.MongoClient = client
	} else {
		log.Fatal(err.Error())
	}
	config.Database = client.Database(database)
	// ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, oidc_provider)
	oidcConfig := &oidc.Config{
		ClientID: oidcClientId,
	}
	verifier := provider.Verifier(oidcConfig)
	config.TokenVerifier = verifier
	return config
}
