package main

import (
	"instawham/controllers"
	"instawham/initializers"
	"instawham/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDb()
	initializers.SyncDB()
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")
	server.Static("/static", "./static")

	server.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"Title": "Hello World"})
	})

	server.GET("/posts", controllers.FetchPosts)

	server.GET("/createpost", func(c *gin.Context) {
		c.HTML(200, "createpost.html", gin.H{"Title": "Create Post"})
	})

	server.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{"Title": "Signup"})
	})
	server.POST("/signup", controllers.SignUp)
	server.POST("/login", controllers.Login)

	server.GET("/validate", middleware.CheckJwt, controllers.Validate)

	server.Run(":8080")
}
