package middleware

import (
	"strings"
	"time"

	"github.com/I000000/recly/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var skipPaths = map[string]bool{
	"/health":  true,
	"/metrics": true,
}

func LoggerMiddleware(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Блокируем спам
		if strings.HasPrefix(path, "/announce") ||
			strings.HasPrefix(path, "/control") ||
			strings.HasPrefix(path, "/checkupdate") ||
			strings.HasPrefix(path, "/scrape") ||
			strings.HasPrefix(path, "/ann") {
			c.AbortWithStatus(403)
			return
		}

		// Пропускаем логирование для определённых путей
		if skipPaths[path] || strings.HasPrefix(path, "/uploads") {
			c.Next()
			return
		}

		// Генерируем request_id
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Добавляем logger с request_id в контекст
		ctx := logger.WithLogger(c.Request.Context(), zapLogger.With(
			zap.String("request_id", requestID),
		))
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Логируем только ошибки или долгие запросы
		if status >= 400 || latency > 1*time.Second {
			zapLogger.Warn("request",
				zap.String("request_id", requestID),
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
			)
		} else if status >= 200 && status < 300 {
			// Логируем успешные запросы только на уровне Debug
			zapLogger.Debug("request",
				zap.String("request_id", requestID),
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("latency", latency),
			)
		}
	}
}
