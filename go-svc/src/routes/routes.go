package routes

import (
	"go-svc/src/controllers"

	"github.com/gin-gonic/gin"
)

// Routes function to serve endpoints
func Routes() {
	route := gin.Default()

	// get all
	route.GET("/api/post", controllers.GetAllPosts)

	// get by id
	route.GET("/api/post/:postId", controllers.GetPostById)

	// create
	route.POST("/api/post", controllers.CreatePost)

	// update
	route.PUT("/api/post/:postId", controllers.UpdatePost)

	// safe delete
	route.DELETE("/api/post/:postId", controllers.DeletePost)

	// Run route whenever triggered
	route.Run()
}
