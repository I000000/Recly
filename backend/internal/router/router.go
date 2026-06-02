package router

import (
	"github.com/I000000/recly/internal/handler"
	"github.com/I000000/recly/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(
	authH *handler.AuthHandler,
	libH *handler.LibraryHandler,
	recH *handler.RecommendationHandler,
	userH *handler.UserHandler,
	searchH *handler.SearchHandler,
	savedItemH *handler.SavedItemHandler,
	secret string,
) *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerWithoutSpam())
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(secret))
		{
			// библиотека
			protected.POST("/book/:id/like", libH.AddBook)
			protected.DELETE("/book/:id/like", libH.RemoveBook)
			protected.GET("/user/library/books", libH.GetBooks)
			protected.POST("/movie/:id/like", libH.AddMovie)
			protected.DELETE("/movie/:id/like", libH.RemoveMovie)
			protected.GET("/user/library/movies", libH.GetMovies)
			protected.GET("/search", searchH.Search)
			protected.GET("/items/batch", searchH.BatchGetItems)

			// рекомендации
			protected.POST("/recommend", recH.Request)
			protected.GET("/user/recommendations/history", recH.GetHistory)
			protected.GET("/result/:taskId", recH.GetResult)

			// закладки
			protected.POST("/user/saved-items", savedItemH.Save)
			protected.DELETE("/user/saved-items/:id", savedItemH.Delete)
			protected.GET("/user/saved-items", savedItemH.Get)

			// профиль
			protected.GET("/user/profile", userH.Profile)
		}
	}
	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
