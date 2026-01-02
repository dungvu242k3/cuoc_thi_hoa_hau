package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cuoc_thi_hoa_hau/internal/adapter/graph/generated"
	"cuoc_thi_hoa_hau/internal/adapter/graph/resolver"
	"cuoc_thi_hoa_hau/internal/adapter/storage/mongodb"
	"cuoc_thi_hoa_hau/internal/core/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// 1. Database Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Replace with your actual connection string
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("beauty_contest")

	// 2. Init Repositories
	contestantRepo := mongodb.NewContestantRepo(db)

	// 3. Init Services
	contestantSvc := service.NewContestantService(contestantRepo)

	scoreRepo := mongodb.NewScoreRepo(db)
	scoringSvc := service.NewScoringService(scoreRepo, contestantRepo)

	// 4. Init Resolver
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		ContestantService: contestantSvc,
		ScoringService:    scoringSvc,
	}}))

	// 5. Router Setup
	r := chi.NewRouter()
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
