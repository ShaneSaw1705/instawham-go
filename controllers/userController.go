package controllers

import (
	"fmt"
	"instawham/initializers"
	"instawham/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get details
	var body struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		fmt.Println("Error reading body", err)
		c.JSON(http.StatusTeapot, gin.H{"Error": err})

		return
	}
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		fmt.Println("Error hashing password", err)
		c.JSON(http.StatusFailedDependency, gin.H{"Error": "Error hashing password"})
		return
	}
	// Look up user
	var previousUser models.User
	initializers.DB.First(&previousUser, "email = ?", body.Email)
	if previousUser.ID != 0 {
		// c.JSON(http.StatusSeeOther, gin.H{"message": "A user already exists under that email"})
		c.HTML(200, "error", gin.H{"message": "A User already exists under that email"})
		return
	}
	// Create user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create user"})
		return
	}
	// Respond
	c.HTML(200, "error", gin.H{"message": "User created! please proceed to login"})
}

func Login(c *gin.Context) {
	// get email and password
	var body struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		fmt.Println("Error reading body", err)
		c.JSON(http.StatusTeapot, gin.H{"Error": err})
		return
	}

	// Look up user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		fmt.Println("Error: invalid email")
		c.JSON(http.StatusTeapot, gin.H{"Error": "invalid email"})
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusFailedDependency, gin.H{"Error": "password is incorrect"})
		return
	}

	// jwt token logic
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(408, gin.H{
			"Error": "failed to create jwt token",
		})
		return
	}

	// Return jwt token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth", tokenString, 3600*24, "", "", false, true)

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusNoContent)
}

func Validate(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, "user unable to be retrieved")
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}
