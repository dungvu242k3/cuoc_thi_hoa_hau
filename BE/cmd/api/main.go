package main

import (
	"context"
	"cuoc_thi_hoa_hau/internal/adapter/cache"
	"cuoc_thi_hoa_hau/internal/adapter/graph/generated"
	mw "cuoc_thi_hoa_hau/internal/adapter/graph/middleware"
	"cuoc_thi_hoa_hau/internal/adapter/graph/resolver"
	filesHandler "cuoc_thi_hoa_hau/internal/adapter/handler"
	"cuoc_thi_hoa_hau/internal/adapter/infra/security"
	"cuoc_thi_hoa_hau/internal/adapter/storage/local"
	"cuoc_thi_hoa_hau/internal/adapter/storage/mongodb"
	"cuoc_thi_hoa_hau/internal/core/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 1. Config & Dependencies
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-me"
	}

	// 2. Database
	dbClient, err := mongodb.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := dbClient.Database("beauty_contest")

	redisClient := cache.NewRedisAdapter(redisAddr, redisPassword, 0)

	// 3. Repos & Services
	// 3. Repos & Services
	contestantRepo := mongodb.NewContestantRepo(db)
	userRepo := mongodb.NewUserRepo(db)
	feedbackRepo := mongodb.NewFeedbackRepo(db)
	scoreRepo := mongodb.NewScoreRepo(db) // Init Score Repo

	// Security Adapters
	hasher := security.NewBcryptHasher()
	tokenProvider := security.NewJWTProvider(jwtSecret)

	contestantSvc := service.NewContestantService(contestantRepo)
	authSvc := service.NewAuthService(userRepo, hasher, tokenProvider)
	feedbackSvc := service.NewFeedbackService(feedbackRepo)
	scoreSvc := service.NewScoringService(scoreRepo, contestantRepo) // Init Score Service

	// 4. GraphQL Resolver
	srvResolver := &resolver.Resolver{
		ContestantSvc: contestantSvc,
		AuthSvc:       authSvc,
		FeedbackSvc:   feedbackSvc,
		ScoreSvc:      scoreSvc, // Inject Score Service
		CacheSvc:      redisClient,
	}

	// 5. Server Setup
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: srvResolver}))

	// Production hardening: Complexity Limit
	srv.Use(extension.FixedComplexityLimit(200))

	// Production hardening: Error Presenter to mask internal errors
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		// Mask internal server errors but allow "public" errors
		// For now, we trust the defaults but you can filter here.
		// Example: if err.Extensions["code"] == "INTERNAL_SERVER_ERROR" { err.Message = "Internal Server Error" }
		return err
	})

	r := chi.NewRouter()

	// CORS configuration (MUST be first to handle OPTIONS requests)
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || allowedOrigins[0] == "" {
		allowedOrigins = []string{"http://localhost:5173", "http://localhost:3000", "https://*", "http://*"}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.Middleware(tokenProvider))
	r.Use(mw.RateLimitMiddleware(redisClient))

	// 6. File Storage Setup (Production Ready: Interface-based)
	// Robust Path Resolution: Find "FE" folder relative to "BE" root
	cwd, _ := os.Getwd()
	uploadDir := filepath.Join(cwd, "..", "FE", "public", "uploads")

	// Helper to find project root if running from cmd/api
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		// Try going up one more level (case: running from cmd/api)
		uploadDir = filepath.Join(cwd, "..", "..", "FE", "public", "uploads")
	}

	// Create directory if not exists (Print log for debugging)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Warning: Failed to create upload dir at %s: %v", uploadDir, err)
	} else {
		log.Printf("File Storage Initialized. Uploads will be saved to: %s", uploadDir)
	}

	localStorage := local.NewLocalStorage(uploadDir, "/uploads")
	fileHandler := filesHandler.NewFileHandler(localStorage)

	// Register Upload Handler
	r.Post("/upload", fileHandler.UploadFile)

	// Serve Static Files (Images)
	// This makes http://localhost:8080/uploads/xxx.jpg reachable
	fs := http.FileServer(http.Dir(uploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
