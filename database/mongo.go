package database

import (
	"context"
	"log"
	models "minna-style-hub/model"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client
var client *mongo.Client

var (
	databaseName   = "mydatabase"
	collectionName = "items"
)

// ConnectToMongoDB connects to MongoDB
func ConnectToMongoDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not found in .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	return nil
}

// GetAllItems retrieves all items from the database
func GetAllItems() ([]models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var items []models.Item

	collection := client.Database(databaseName).Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// AddItem adds a new item to the database
func AddItem(item models.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(databaseName).Collection(collectionName)

	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

// UpdateItem updates an existing item in the database
// UpdateItem updates an existing item in the database
func UpdateItem(item models.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(databaseName).Collection(collectionName)

	// Exclude _id field from update
	update := bson.M{
		"$set": bson.M{
			"title":      item.Title,
			"text":       item.Text,
			"brand":      item.Brand,
			"images":     item.Images,
			"buttonLink": item.ButtonLink,
			// Add other fields you want to update here
		},
	}

	// Use item.ID directly
	filter := bson.M{"_id": item.ID}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteItem deletes an item from the database by its MongoDB _id
func DeleteItem(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(databaseName).Collection(collectionName)

	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func GetItem(id string) (models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(databaseName).Collection(collectionName)

	var item models.Item

	filter := bson.M{"_id": id}
	err := collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle case where document is not found
			return models.Item{}, err
		}
		return models.Item{}, err
	}

	return item, nil
}
