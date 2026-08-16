package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/I000000/recly/internal/config"
	"github.com/I000000/recly/internal/handler"
	"github.com/I000000/recly/internal/meili"
	"github.com/I000000/recly/internal/rabbitmq"
	redisPkg "github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/internal/repository/postgres"
	"github.com/I000000/recly/internal/router"
	"github.com/I000000/recly/internal/service"
)

// @title Recly API
// @version 1.0
// @description Cross-media recommendation system API

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// --- Настройка логгера с ротацией ---
	logDir := filepath.Dir("logs/app.log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	logWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,  // 100 MB
		MaxBackups: 5,    // хранить 5 старых файлов
		MaxAge:     30,   // хранить 30 дней
		Compress:   true, // сжимать старые логи
	}
	defer logWriter.Close()

	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout", logWriter.Filename}
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("application starting", zap.String("version", "1.0.0"))

	// --- Подключения ---
	pool, err := postgres.NewPool(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("database", zap.Error(err))
	}
	defer pool.Close()

	publisher, err := rabbitmq.NewAMQPPublisher(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatal("rabbitmq", zap.Error(err))
	}
	defer publisher.Close()

	cache := redisPkg.NewRedisCache(cfg.RedisURL, "", 0)

	meiliClient := meili.NewClient("http://meilisearch:7700", "aSecretMasterKey")

	// --- Репозитории ---
	userRepo := postgres.NewUserRepo(pool)
	tokenRepo := postgres.NewTokenRepo(pool)
	libRepo := postgres.NewLibraryRepo(pool)
	recRepo := postgres.NewRecommendationRepo(pool)
	savedItemRepo := postgres.NewSavedItemRepo(pool)
	viewedItemRepo := postgres.NewViewedItemRepo(pool)

	// --- Сервисы ---
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	libSvc := service.NewLibraryService(libRepo)
	sourceSelectorSvc := service.NewSourceSelector(libSvc)
	recSvc := service.NewRecommendationService(recRepo, publisher, cache, sourceSelectorSvc)
	searchSvc := service.NewSearchService(meiliClient)
	savedItemSvc := service.NewSavedItemService(savedItemRepo)
	viewedItemSvc := service.NewViewedItemService(viewedItemRepo)

	// --- Хендлеры ---
	authH := handler.NewAuthHandler(authSvc)
	libH := handler.NewLibraryHandler(libSvc)
	recH := handler.NewRecommendationHandler(recSvc)
	userH := handler.NewUserHandler(userSvc)
	searchH := handler.NewSearchHandler(searchSvc)
	savedItemH := handler.NewSavedItemHandler(savedItemSvc)
	viewedItemH := handler.NewViewedItemHandler(viewedItemSvc)

	// --- Роутер ---
	r := router.Setup(authH, libH, recH, userH, searchH, savedItemH, viewedItemH, cfg.JWTSecret, logger)
	r.Static("/uploads", "./uploads")

	srv := &http.Server{Addr: ":" + cfg.ServerPort, Handler: r}

	go func() {
		logger.Info("Server starting on " + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen", zap.Error(err))
		}
	}()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown", zap.Error(err))
	}
	logger.Info("Server stopped")
}
