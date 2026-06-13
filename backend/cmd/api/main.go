package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/I000000/recly/internal/config"
	"github.com/I000000/recly/internal/handler"
	"github.com/I000000/recly/internal/rabbitmq"
	redisPkg "github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/internal/repository/postgres"
	"github.com/I000000/recly/internal/router"
	"github.com/I000000/recly/internal/service"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	pool, err := postgres.NewPool(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("database", zap.Error(err))
	}
	defer pool.Close()

	// RabbitMQ Publisher
	publisher, err := rabbitmq.NewAMQPPublisher(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatal("rabbitmq", zap.Error(err))
	}
	defer publisher.Close()

	// Redis Cache
	cache := redisPkg.NewRedisCache(cfg.RedisURL, "", 0)

	// Репозитории
	userRepo := postgres.NewUserRepo(pool)
	tokenRepo := postgres.NewTokenRepo(pool)
	libRepo := postgres.NewLibraryRepo(pool)
	recRepo := postgres.NewRecommendationRepo(pool)
	savedItemRepo := postgres.NewSavedItemRepo(pool)
	viewedItemRepo := postgres.NewViewedItemRepo(pool)

	// Сервисы
	userService := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	libSvc := service.NewLibraryService(libRepo)
	recSvc := service.NewRecommendationService(recRepo, publisher, cache, libSvc)
	searchSvc := service.NewSearchService("http://meilisearch:7700", "aSecretMasterKey")
	savedItemSvc := service.NewSavedItemService(savedItemRepo)
	viewedItemSvc := service.NewViewedItemService(viewedItemRepo)

	// Хэндлеры
	authH := handler.NewAuthHandler(authSvc)
	libH := handler.NewLibraryHandler(libSvc)
	recH := handler.NewRecommendationHandler(recSvc)
	userH := handler.NewUserHandler(userService)
	searchH := handler.NewSearchHandler(searchSvc)
	savedItemH := handler.NewSavedItemHandler(savedItemSvc)
	viewedItemH := handler.NewViewedItemHandler(viewedItemSvc)

	r := router.Setup(authH, libH, recH, userH, searchH, savedItemH, viewedItemH, cfg.JWTSecret)

	r.Static("/uploads", "./uploads")

	srv := &http.Server{Addr: ":" + cfg.ServerPort, Handler: r}
	go func() {
		logger.Info("Server starting on " + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen", zap.Error(err))
		}
	}()

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
