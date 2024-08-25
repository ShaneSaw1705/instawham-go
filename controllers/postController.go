package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func FetchPosts(c *gin.Context) {
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
}
