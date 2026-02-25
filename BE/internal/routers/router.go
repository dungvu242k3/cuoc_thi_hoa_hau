/*
Package routers chịu trách nhiệm Trạm điều phối giao thông (Routing).
Tác dụng: Khi một thiết bị gửi request tới url ("/api/auth/login"), Router sẽ phân theo đúng map
để điều hướng request đó vào đúng Handler (Controller) tương ứng.
Đồng thời, Router cung cấp Middleware (Trạm gác) như CORS, Logging, hoặc Auth chặn trước cửa.
*/
package routers

import (
	"net/http"
	"os"
	"path/filepath"

	mw "cuoc_thi_hoa_hau/graph/middleware"
	"cuoc_thi_hoa_hau/internal/config"
	handlers "cuoc_thi_hoa_hau/internal/handler"
	"cuoc_thi_hoa_hau/internal/types"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Deps struct {
	GraphQLHandler http.Handler
	AuthHandler    *handlers.AuthHandler
	FileHandler    *handlers.FileHandler
	TokenProvider  types.TokenProvider
	CacheSvc       types.CacheService
	CORS           config.CORSConfig
	UploadDir      string
}

func RegisterRoutes(r *chi.Mux, deps *Deps) {
	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.CORS.AllowedOrigins,
		AllowedMethods:   deps.CORS.AllowedMethods,
		AllowedHeaders:   deps.CORS.AllowedHeaders,
		AllowCredentials: deps.CORS.AllowCredentials,
		MaxAge:           deps.CORS.MaxAge,
	}))

	// Auth & client info middleware
	r.Use(mw.Middleware(deps.TokenProvider))
	r.Use(mw.ClientInfoMiddleware)
	if deps.CacheSvc != nil {
		r.Use(mw.RateLimitMiddleware(deps.CacheSvc))
	}

	// GraphQL
	r.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", deps.GraphQLHandler)

	// REST API: Mapping các URL cụ thể -> Hàm xử lý (Handler)
	r.Post("/api/auth/register", deps.AuthHandler.RegisterAudience)
	r.Post("/api/auth/login", deps.AuthHandler.LoginAudience)
	r.Post("/api/upload", deps.FileHandler.UploadFile)

	// Static files: Cho phép duyệt file public, ví dụ ảnh đại diện
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, deps.UploadDir)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(filesDir))))
}
