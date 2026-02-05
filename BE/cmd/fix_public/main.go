package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env
	_ = godotenv.Load("../../.env")

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("beauty_contest")
	collContestant := db.Collection("contestants")

	// Target the specific user ID or just update all pending to approved+public for testing
	// Based on screenshot, ID is 695e58beb5c7ae6f869739b8
	contestantID := "695e58beb5c7ae6f869739b8"

	oid, _ := primitive.ObjectIDFromHex(contestantID)

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"is_public": true,
			"status":    "approved", // Use string directly to match DB
		},
	}

	res, err := collContestant.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("Failed to update contestant: %v", err)
	}

	fmt.Printf("Matched %d documents and modified %d documents.\n", res.MatchedCount, res.ModifiedCount)
	if res.ModifiedCount > 0 {
		fmt.Println("SUCCESS: Contestant is now APPROVED and PUBLIC.")
	} else {
		fmt.Println("WARNING: No document was modified (maybe ID is wrong or already updated).")
	}
}
