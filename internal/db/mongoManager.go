package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/erdemkosk/box_bridge/internal/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoManager struct {
	client   *mongo.Client
	dbName   string
	database *mongo.Database
}

func NewMongoManager(mongoDbURL string) (*MongoManager, error) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoDbURL)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("BOX-BRIDGE: MongoDB connection failed: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("BOX-BRIDGE: Unable to establish MongoDB connection: %v", err)
	}

	dbName := getDatabaseName(mongoDbURL)

	database := client.Database(dbName)

	if err := createCollectionIfNotExists(ctx, database, "inbox"); err != nil {
		log.Printf("BOX-BRIDGE: Error occurred while creating inbox collection: %v", err)
	}
	if err := createCollectionIfNotExists(ctx, database, "outbox"); err != nil {
		log.Printf("BOX-BRIDGE: Error occurred while creating outbox collection: %v", err)
	}

	log.Println("BOX-BRIDGE: MongoDB connection successfully established.")

	return &MongoManager{
		client:   client,
		dbName:   dbName,
		database: database,
	}, nil
}

func (m *MongoManager) insertMessage(collName string, message interface{}) error {
	collection := m.database.Collection(collName)

	_, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		return fmt.Errorf("BOX-BRIDGE: %s message insertion error: %v", collName, err)
	}

	return nil
}

func (m *MongoManager) SaveToOutbox(message models.Outbox) error {
	return m.insertMessage("outbox", message)
}

func (m *MongoManager) SaveToInbox(message models.Inbox) error {
	return m.insertMessage("inbox", message)
}

func (m *MongoManager) UpdateOutboxStatus(messageID string, status string) error {
	collection := m.database.Collection("outbox")

	filter := bson.M{"message_id": messageID}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("BOX-BRIDGE: failed to update message status: %v", err)
	}

	return nil
}

func createCollectionIfNotExists(ctx context.Context, db *mongo.Database, collName string) error {
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("BOX-BRIDGE: could not retrieve collections: %v", err)
	}

	for _, collection := range collections {
		if collection == collName {
			return nil
		}
	}

	err = db.CreateCollection(ctx, collName)
	if err != nil {
		return fmt.Errorf("BOX-BRIDGE: could not create collection: %v", err)
	}

	log.Printf("BOX-BRIDGE: Collection '%s' created successfully", collName)
	return nil
}

func getDatabaseName(connectionString string) string {
	parsedURL, err := url.Parse(connectionString)
	if err != nil {
		log.Fatalf("BOX-BRIDGE: connection string parse problem: %w", err)
	}

	if parsedURL.Path == "" || parsedURL.Path == "/" {
		log.Fatalf("BOX-BRIDGE: cannot find db name")
	}

	return parsedURL.Path[1:]
}

func (m *MongoManager) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("BOX-BRIDGE: MongoDB disconnection failed: %v", err)
	}

	log.Println("BOX-BRIDGE: MongoDB connection successfully closed.")
	return nil
}
