package routes

import (
	"github.com/gin-gonic/gin"
	"post-comments-api/controllers"
	"post-comments-api/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.LoggerMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)

		api.GET("/users/me", middleware.AuthMiddleware(), controllers.GetCurrentUser)

		// Public posts/comments
		public := api.Group("/public")
		public.POST("/posts", controllers.CreatePostPublic)
		public.POST("/comments", controllers.CreateCommentPublic)

		// Protected posts
		api.GET("/posts", controllers.GetPosts)
		api.POST("/posts", middleware.AuthMiddleware(), controllers.CreatePost)
		api.GET("/posts/:id", controllers.GetPost)
		api.PUT("/posts/:id", middleware.AuthMiddleware(), controllers.UpdatePost)
		api.DELETE("/posts/:id", middleware.AuthMiddleware(), controllers.DeletePost)

		// Comments
		api.GET("/posts/:id/comments", controllers.GetComments)
		api.POST("/posts/:id/comments", middleware.AuthMiddleware(), controllers.CreateComment)
		api.POST("/comments", middleware.AuthMiddleware(), controllers.CreateComment)
		api.PUT("/comments/:id", middleware.AuthMiddleware(), controllers.UpdateComment)
		api.DELETE("/comments/:id", middleware.AuthMiddleware(), controllers.DeleteComment)
	}

	return r
}
