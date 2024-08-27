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
	responseChan := make(chan []models.Post)

	go func() {
		var posts []models.Post
		result := initializers.DB.Order("created_at desc").Limit(50).Find(&posts)
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

func FetchSinglePost(c *gin.Context) {
	responseChan := make(chan models.Post)
	id := c.Param("id")

	go func() {
		var post models.Post
		result := initializers.DB.First(&post, "ID = ?", id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post no longer exists"})
			return
		}
		responseChan <- post
	}()

	post := <-responseChan
	c.JSON(200, post)
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

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusNoContent)
}
