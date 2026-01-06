package main

import (
	"go-note-editor/controllers"
	"go-note-editor/initializers"
	"go-note-editor/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()

	// Public auth routes
	r.POST("/auth/signup", controllers.Signup)
	r.POST("/auth/signin", controllers.Signin)

	// Protected routes (require authentication)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/notes", controllers.GetAllNotes)
		protected.POST("/notes", controllers.CreateNote)
		protected.GET("/notes/:id", controllers.GetNoteByID)
		protected.PUT("/notes/:id", controllers.UpdateNote)
		protected.DELETE("/notes/:id", controllers.DeleteNote)
	}

	r.Run(":8080")
}
