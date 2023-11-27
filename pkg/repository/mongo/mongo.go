package mongo

import (
	"context"
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"medbuddy-backend/internal/config"
	"medbuddy-backend/pkg/repository/storage"
	"medbuddy-backend/utility"
	"time"
)

var (
	mongoclient         *mongo.Client
	generalQueryTimeout = 60 * time.Second

	logger = utility.NewLogger()
)

type Mongo struct {
	mongoclient *mongo.Client
	timeout     time.Duration
}

func GetDB() storage.StorageRepository {
	return &Mongo{mongoclient: mongoclient, timeout: generalQueryTimeout}
}

func Connection() (db *mongo.Client) {
	return mongoclient
}

func ConnectToDB() *mongo.Client {
	ctx := context.Background()

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	uri := config.GetConfig().MongoHost
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri), options.Client().SetTLSConfig(tlsConfig))
	if err != nil {
		log.Fatal("Error connecting to mongoDB, error: ", err)
	}

	// PINGING THE CONNECTION
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error pinging mongoDB connection, error: ", err)
	}

	// IF EVERYTHING IS OKAY, THEN CONNECTION IS SETTLED
	logger.Info("MONGO CONNECTION ESTABLISHED")

	mongoclient = mongoClient
	return mongoClient
}

func DisconnectDB(ctx context.Context) {
	err := mongoclient.Disconnect(ctx)
	if err != nil {
		logger.Error("Error disconnecting DB connection, error: ", err)
		return
	}
	logger.Info("SUCCESSFULLY CLOSED DB CONNECTION")
}
