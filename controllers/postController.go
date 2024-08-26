package controllers

import (
	"instawham/initializers"
	"instawham/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func FetchPost(c *gin.Context) {
	responseChan := make(chan []Post)

	go func() {
		var posts []Post
		result := initializers.DB.Find(&posts)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "failed to find posts"})
		}

		responseChan <- posts
	}()
	posts := <-responseChan
	if posts == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "failed to read posts"})
	}

	c.HTML(200, "posts", gin.H{"posts": posts})
}

func CreatePost(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "user not found"})
		return
	}

	authenticatedUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "user type assertion failed"})
		return
	}

	var body struct {
		Title       string `form:"title"`
		Description string `form:"description"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error reading body"})
		return
	}

	post := models.Post{
		Title:       body.Title,
		Description: body.Description,
		AuthorID:    int(authenticatedUser.ID),
	}

	initializers.DB.Create(&post)

	c.JSON(http.StatusOK, gin.H{"Success": "post created"})
}
