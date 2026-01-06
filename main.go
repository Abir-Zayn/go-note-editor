package main

import (
	"go-note-editor/controllers"
	"go-note-editor/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()

	// Auth routes
	r.POST("/auth/signup", controllers.Signup)

	// Note routes (existing)
	r.GET("/notes", controllers.GetAllNotes)
	r.POST("/notes", controllers.CreateNote)

	r.Run(":8080")
}
