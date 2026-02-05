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
	collUser := db.Collection("users")
	collContestant := db.Collection("contestants")

	fmt.Println("=== DIAGNOSTIC REPORT ===")

	// 1. Check Users
	totalUsers, _ := collUser.CountDocuments(ctx, bson.M{})
	fmt.Printf("Total Users: %d\n", totalUsers)

	// List ALL Users to check Roles
	cursor, _ := collUser.Find(ctx, bson.M{})
	var users []bson.M
	cursor.All(ctx, &users)

	fmt.Printf("Listing all %d users:\n", len(users))
	for _, u := range users {
		uid := u["_id"]
		username := u["username"]
		role := u["role_id"]
		fmt.Printf(" - UserID: %v | Username: %v | Role: %v\n", uid, username, role)

		// Check linked profile if it looks like a candidate/contestant
		if role == "candidate" || role == "contestant" {
			var profile bson.M
			err := collContestant.FindOne(ctx, bson.M{"user_id": uid.(primitive.ObjectID).Hex()}).Decode(&profile)
			if err == nil {
				isPublic := profile["is_public"]
				status := profile["status"]
				fmt.Printf("   -> HAS PROFILE. Status: %v, IsPublic: %v\n", status, isPublic)
			} else {
				fmt.Println("   -> [WARNING] NO CONTESTANT PROFILE FOUND.")
			}
		}
	}

	// Check for any public contestants
	publicCount, _ := collContestant.CountDocuments(ctx, bson.M{"is_public": true})
	fmt.Printf("\nTotal Public Contestants (Visible on Dashboard): %d\n", publicCount)

	fmt.Println("=========================")
}
