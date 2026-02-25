/*
Package config chịu trách nhiệm quản lý cấu hình tập trung cho toàn nền tảng.
Tác dụng:
  - Thay vì mọi file đều gọi os.Getenv("MONGO_URI") phân tán khắp nơi, hệ thống sẽ gom tất cả cấu hình
    vào 1 Struct chuẩn (Config).
  - Mọi hàm cần dùng cấu hình chỉ việc gọi cfg.Database.MongoDB.DSN cực kì an toàn và có gợi ý code (typing).
*/
package config

import (
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	HTTP     HTTPConfig     `yaml:"http"`
	CORS     CORSConfig     `yaml:"cors"`
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
	Host    string `yaml:"host"`
}

type HTTPConfig struct {
	Port    int `yaml:"port"`
	Timeout int `yaml:"timeout"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowedOrigins"`
	AllowedMethods   []string `yaml:"allowedMethods"`
	AllowedHeaders   []string `yaml:"allowedHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	IsSave bool   `yaml:"isSave"`
}

type DatabaseConfig struct {
	Driver  string      `yaml:"driver"`
	MongoDB MongoConfig `yaml:"mongodb"`
}

type MongoConfig struct {
	DSN    string `yaml:"dsn"`
	DBName string `yaml:"dbName"`
}

type RedisConfig struct {
	DSN          string `yaml:"dsn"`
	DialTimeout  int    `yaml:"dialTimeout"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expireHours"`
}

type UploadConfig struct {
	Dir          string   `yaml:"dir"`
	PublicURL    string   `yaml:"publicURL"`
	MaxSizeMB    int      `yaml:"maxSizeMB"`
	AllowedTypes []string `yaml:"allowedTypes"`
}

// Load đọc file configs/passiontech-beauty-contest.yml.
// Trong yaml có sử dụng ${BÍ_DANH}, Go sẽ thay thế nó bằng giá trị thực tế của biến môi trường (os.ExpandEnv).
func Load() *Config {
	cfg := &Config{}

	data, err := os.ReadFile("configs/passiontech-beauty-contest.yml")
	if err != nil {
		log.Printf("[Config] YAML not found, falling back to env vars: %v", err)
		return loadFromEnv()
	}

	// Expand environment variables in YAML (e.g. ${MONGO_URI})
	expanded := os.ExpandEnv(string(data))

	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		log.Fatalf("[Config] Failed to parse YAML: %v", err)
	}

	// Defaults
	if cfg.HTTP.Port == 0 {
		cfg.HTTP.Port = 8080
	}
	if cfg.Upload.Dir == "" {
		cfg.Upload.Dir = "./uploads"
	}
	if cfg.Upload.PublicURL == "" {
		cfg.Upload.PublicURL = "http://localhost:8080/uploads"
	}
	if cfg.JWT.ExpireHours == 0 {
		cfg.JWT.ExpireHours = 72
	}

	return cfg
}

// loadFromEnv fallback: load config from environment variables (backward compatible)
func loadFromEnv() *Config {
	port := 8080
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	publicURL := os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:8080/uploads"
	}

	return &Config{
		App:  AppConfig{Name: "passiontech-beauty-contest", Env: "dev"},
		HTTP: HTTPConfig{Port: port},
		CORS: CORSConfig{
			AllowedOrigins:   strings.Split(os.Getenv("CORS_ORIGINS"), ","),
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		},
		Database: DatabaseConfig{
			Driver: "mongodb",
			MongoDB: MongoConfig{
				DSN:    os.Getenv("MONGO_URI"),
				DBName: os.Getenv("MONGO_DB"),
			},
		},
		Redis: RedisConfig{DSN: os.Getenv("REDIS_ADDR")},
		JWT:   JWTConfig{Secret: os.Getenv("JWT_SECRET"), ExpireHours: 72},
		Upload: UploadConfig{
			Dir:       uploadDir,
			PublicURL: publicURL,
			MaxSizeMB: 10,
		},
	}
}
