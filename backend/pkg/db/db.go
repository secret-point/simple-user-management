// File: /backend/pkg/db/db.go

package db

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConnection holds the database connection instance.
var DBConnection *mongo.Client

// ConnectDB initializes the connection to MongoDB.
func ConnectDB() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://secretpoint2020:Ibegyourpardon@cluster0.7uwlvlg.mongodb.net"))
	if err != nil {
			log.Error("Failed to connect to database:", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
			log.Error("Failed to ping database:", err)
	}

	log.Info("Successfully connected to MongoDB.")
	return client
}
