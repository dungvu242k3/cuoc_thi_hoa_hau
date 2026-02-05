package main

import (
	"context"
	"cuoc_thi_hoa_hau/internal/adapter/infra/security"
	"cuoc_thi_hoa_hau/internal/core/domain"
	"encoding/json"
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
	_ = godotenv.Load("../../.env") // Try loading from root if running from cmd/seeder

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
	collRole := db.Collection("roles")
	collPerm := db.Collection("permissions")

	hasher := security.NewBcryptHasher()

	// 1. Seed Permissions
	perms := []domain.Permission{
		{Code: domain.PermScoreWrite, Description: "Chấm điểm thí sinh"},
		{Code: domain.PermScoreRead, Description: "Xem điểm"},
		{Code: domain.PermContestantRead, Description: "Xem hồ sơ thí sinh"},
		{Code: domain.PermContestantWrite, Description: "Quản lý hồ sơ thí sinh"},
		{Code: domain.PermSystemConfig, Description: "Cấu hình hệ thống"},
		{Code: domain.PermUserRead, Description: "Xem danh sách người dùng"},
	}

	for _, p := range perms {
		_, err := collPerm.ReplaceOne(ctx, bson.M{"_id": p.Code}, p, options.Replace().SetUpsert(true))
		if err != nil {
			log.Printf("Failed to seed perm %s: %v", p.Code, err)
		}
	}
	fmt.Println("Seeded Permissions.")

	// 2. Seed Roles
	roles := []domain.RoleDef{
		{
			ID:   "admin",
			Name: "Quản trị viên",
			Permissions: []string{
				domain.PermScoreRead, domain.PermContestantRead, domain.PermContestantWrite,
				domain.PermSystemConfig, domain.PermUserRead,
			},
		},
		{
			ID:   "examiner",
			Name: "Ban Giám Khảo",
			Permissions: []string{
				domain.PermScoreWrite, domain.PermScoreRead, domain.PermContestantRead,
			},
		},
		{
			ID:   "candidate",
			Name: "Thí sinh",
			Permissions: []string{
				domain.PermContestantRead,
			},
		},
	}

	for _, r := range roles {
		r.CreatedAt = time.Now()
		r.UpdatedAt = time.Now()
		_, err := collRole.ReplaceOne(ctx, bson.M{"_id": r.ID}, r, options.Replace().SetUpsert(true))
		if err != nil {
			log.Printf("Failed to seed role %s: %v", r.ID, err)
		}
	}
	fmt.Println("Seeded Roles.")

	// 3. Seed Examiners (From Env)
	type ExaminerSeed struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var examiners []ExaminerSeed

	seedData := os.Getenv("SEED_EXAMINERS")
	if seedData != "" {
		if err := json.Unmarshal([]byte(seedData), &examiners); err != nil {
			log.Printf("Failed to parse SEED_EXAMINERS: %v", err)
		}
	} else {
		log.Println("No SEED_EXAMINERS env found, skipping examiner seeding.")
	}

	for _, ex := range examiners {
		// Check exists
		count, _ := collUser.CountDocuments(ctx, bson.M{"username": ex.Username})
		if count > 0 {
			log.Printf("Examiner %s already exists, skipping...", ex.Username)
			continue
		}

		hashed, _ := hasher.Hash(ex.Password)
		user := domain.User{
			Username:  ex.Username,
			Password:  hashed,
			RoleID:    "examiner", // domain.RoleExaminer
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := collUser.InsertOne(ctx, user)
		if err != nil {
			log.Printf("Failed to create %s: %v", ex.Username, err)
		} else {
			fmt.Printf("Created examiner: %s (pass: %s)\n", ex.Username, "******")
		}
	}

	// Also ensure an admin exists
	if count, _ := collUser.CountDocuments(ctx, bson.M{"username": "admin@event.com"}); count == 0 {
		hashed, _ := hasher.Hash("Admin@123")
		collUser.InsertOne(ctx, domain.User{
			Username:  "admin@event.com",
			Password:  hashed,
			RoleID:    "admin",
			CreatedAt: time.Now(),
		})
		fmt.Println("Created admin: admin@event.com (pass: Admin@123)")
	}

	// 4. Seed Contestants (Candidates)
	contestantCount := 10
	fmt.Printf("Seeding %d contestants...\n", contestantCount)

	// We need the Contestant Collection
	collContestant := db.Collection("contestants")

	for i := 1; i <= contestantCount; i++ {
		email := fmt.Sprintf("thiSinh%d@event.com", i)
		sbd := fmt.Sprintf("%03d", i)

		// 4.1 Create User
		var userID string
		var existingUser domain.User
		err := collUser.FindOne(ctx, bson.M{"username": email}).Decode(&existingUser)
		if err == nil {
			userID = existingUser.ID
			// log.Printf("User %s exists", email)
		} else {
			hashed, _ := hasher.Hash("123456")
			newUser := domain.User{
				Username:  email,
				Password:  hashed,
				RoleID:    "candidate",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			res, err := collUser.InsertOne(ctx, newUser)
			if err != nil {
				log.Printf("Failed to create user %s: %v", email, err)
				continue
			}
			userID = res.InsertedID.(primitive.ObjectID).Hex()
		}

		// 4.2 Create Contestant Profile
		count, _ := collContestant.CountDocuments(ctx, bson.M{"user_id": userID})
		if count > 0 {
			continue
		}

		profile := domain.Contestant{
			UserID:   userID,
			SBD:      sbd,
			Status:   domain.ContestantStatusApproved,
			IsPublic: true,
			PersonalInfo: domain.PersonalInfo{
				FullName:    fmt.Sprintf("Nguyễn Thị Thí Sinh %d", i),
				DateOfBirth: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Nationality: "Việt Nam",
				Email:       email,
				Phone:       fmt.Sprintf("0900000%03d", i),
				Address:     "Hà Nội",
				Job:         "Sinh viên",
				Gender:      "Nữ",
			},
			PhysicalInfo: domain.PhysicalInfo{
				Height:       165 + float64(i),
				Weight:       50 + float64(i),
				Measurements: "90-60-90",
			},
			Portfolio: domain.Portfolio{
				AvatarURL:    fmt.Sprintf("https://avatar.iran.liara.run/public/girl?username=%s", sbd),
				Introduction: "Xin chào, tôi là thí sinh hoa hậu.",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = collContestant.InsertOne(ctx, profile)
		if err != nil {
			log.Printf("Failed to create profile for %s: %v", email, err)
		} else {
			fmt.Printf("Created contestant: %s (SBD: %s)\n", email, sbd)
		}
	}

	fmt.Println("Seeding completed.")
}
