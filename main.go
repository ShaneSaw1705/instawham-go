package main

import (
	"encoding/json"
	"fmt"
	"instawham/controllers"
	"instawham/initializers"
	"instawham/middleware"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

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
		c.HTML(200, "index.html", gin.H{
			"Title": "Hello World",
		})
	})

	server.GET("/posts", func(c *gin.Context) {
		responseChan := make(chan []Post)

		go func() {
			resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
			if err != nil {
				fmt.Println("Error fetching from address", err)
				responseChan <- nil
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("error reading body", err)
				responseChan <- nil
				return
			}

			var posts []Post
			if err := json.Unmarshal(body, &posts); err != nil {
				fmt.Println("Error unmarshalling JSON", err)
				responseChan <- nil
				return
			}

			responseChan <- posts
		}()

		posts := <-responseChan
		if posts == nil {
			c.JSON(500, gin.H{"error": "Error fetching posts"})
			return
		}

		c.HTML(200, "posts", gin.H{"posts": posts})
	})

	server.GET("/createpost", func(c *gin.Context) {
		c.HTML(200, "createpost.html", gin.H{
			"Title": "Create Post",
		})
	})

	server.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{"Title": "Signup"})
	})
	server.POST("/signup", controllers.SignUp)
	server.POST("/login", controllers.Login)

	server.GET("/validate", middleware.CheckJwt, controllers.Validate)

	server.Run(":8080")
}
