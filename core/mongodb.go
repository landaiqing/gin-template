package core

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schisandra-cloud-album/global"
	"time"
)

// InitMongoDB initializes the MongoDB connection
func InitMongoDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(global.CONFIG.MongoDB.MongoDsn())
	clientOptions.SetAuth(options.Credential{
		AuthMechanism:           "SCRAM-SHA-256",
		AuthMechanismProperties: nil,
		AuthSource:              global.CONFIG.MongoDB.AuthSource,
		Username:                global.CONFIG.MongoDB.User,
		Password:                global.CONFIG.MongoDB.Password,
		PasswordSet:             true,
	})
	clientOptions.SetConnectTimeout(time.Duration(global.CONFIG.MongoDB.Timeout) * time.Second)
	clientOptions.SetMaxPoolSize(global.CONFIG.MongoDB.MaxOpenConn)
	clientOptions.SetMaxConnecting(global.CONFIG.MongoDB.MaxIdleConn)
	connect, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		global.LOG.Fatalf(err.Error())
		return
	}
	// Check the connection
	err = connect.Ping(context.TODO(), nil)
	if err != nil {
		global.LOG.Fatalf(err.Error())
		return
	}
	global.MongoDB = connect
}
