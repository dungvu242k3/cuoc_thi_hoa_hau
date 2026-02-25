/*
Package initial đóng vai trò là "Nhà máy lắp ráp" (Dependency Injection Container) của ứng dụng.
Tác dụng: Khởi tạo tất cả các lớp (Layer) từ dưới lên trên và tiêm (inject) chúng vào nhau.
Thứ tự chuẩn:
1. Load Config (Cấu hình)
2. Connect Database / Cache (Hạ tầng)
3. Init DAO / Repository (Tầng kết nối DB)
4. Init Service (Tầng nghiệp vụ, được tiêm DAO vào)
5. Init Handler (Tầng HTTP/GraphQL, được tiêm Service vào)
6. Init Router (Đăng ký URL) và Start Server.
Nó giúp các file code ở bên trong không bị dính chặt (coupled) vào nhau.
*/
package initial

import (
	"context"
	"fmt"
	"log"

	"cuoc_thi_hoa_hau/graph/generated"
	"cuoc_thi_hoa_hau/graph/resolver"
	"cuoc_thi_hoa_hau/internal/cache"
	"cuoc_thi_hoa_hau/internal/config"
	"cuoc_thi_hoa_hau/internal/dao"
	"cuoc_thi_hoa_hau/internal/database"
	handlers "cuoc_thi_hoa_hau/internal/handler"
	"cuoc_thi_hoa_hau/internal/routers"
	"cuoc_thi_hoa_hau/internal/server"
	"cuoc_thi_hoa_hau/internal/service"

	"net/http"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

// InitAndRun initializes all dependencies and starts the HTTP server
func InitAndRun() {
	// 1. Tải cấu hình từ file YAML hoặc .env vào struct Config
	cfg := config.Load()

	// Validate required config
	if cfg.Database.MongoDB.DSN == "" || cfg.Database.MongoDB.DBName == "" {
		log.Fatal("Database MongoDB DSN and DBName must be set")
	}
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT secret must be set")
	}

	// Database
	client, err := database.Connect(cfg.Database.MongoDB.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(cfg.Database.MongoDB.DBName)

	// 3. Tầng Security: Khởi tạo các công cụ băm mật khẩu và tạo token JWT
	hasher := server.NewBcryptHasher()
	tokenProvider := server.NewJWTProvider(cfg.JWT.Secret)

	// 4. Tầng DAO (Data Access Object): Khởi tạo các "Repo" chuyên nói chuyện với CSDL MongoDB
	userRepo := dao.NewUserRepo(db)
	contestantRepo := dao.NewContestantRepo(db)
	feedbackRepo := dao.NewFeedbackRepo(db)
	scoreRepo := dao.NewScoreRepo(db)
	scheduleRepo := dao.NewScheduleRepo(db)

	// 5. Tầng Service (Business Logic): Chứa quy tắc kinh doanh (ví dụ: tuổi > 18 mới tạo thí sinh)
	// Service cần gọi DB, nên ta "tiêm" Repo vào Service
	authSvc := service.NewAuthService(userRepo, hasher, tokenProvider)
	contestantSvc := service.NewContestantService(contestantRepo)
	feedbackSvc := service.NewFeedbackService(feedbackRepo)
	scoreSvc := service.NewScoringService(scoreRepo, contestantRepo)
	scheduleSvc := service.NewScheduleService(scheduleRepo)

	// Cache
	cacheSvc := cache.NewRedisAdapter(cfg.Redis.DSN, "", 0)

	// File Storage
	localStorage := dao.NewLocalStorage(cfg.Upload.Dir, cfg.Upload.PublicURL)

	// 6. Tầng Handler (Controller): Chuyên nhận Request HTTP và gọi Service
	// Handler chỉ lo việc phân tích JSON đầu vào, gọi Service, rồi trả về JSON kết quả
	authHandler := handlers.NewAuthHandler(authSvc)
	fileHandler := handlers.NewFileHandler(localStorage)

	// GraphQL
	graphResolver := &resolver.Resolver{
		ContestantSvc: contestantSvc,
		AuthSvc:       authSvc,
		ScoreSvc:      scoreSvc,
		ScheduleSvc:   scheduleSvc,
		FeedbackSvc:   feedbackSvc,
		CacheSvc:      cacheSvc,
	}
	graphqlHandler := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: graphResolver,
	}))

	// Router
	r := chi.NewRouter()
	routers.RegisterRoutes(r, &routers.Deps{
		GraphQLHandler: graphqlHandler,
		AuthHandler:    authHandler,
		FileHandler:    fileHandler,
		TokenProvider:  tokenProvider,
		CacheSvc:       cacheSvc,
		CORS:           cfg.CORS,
		UploadDir:      cfg.Upload.Dir,
	})

	// Start Server
	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	log.Printf("[%s] %s starting on %s (env=%s)", cfg.App.Name, cfg.App.Version, addr, cfg.App.Env)
	log.Fatal(http.ListenAndServe(addr, r))
}
