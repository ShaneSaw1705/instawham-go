package middleware

import (
	"fmt"
	"instawham/initializers"
	"instawham/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckJwt(c *gin.Context) {
	tokenString, err := c.Cookie("auth")
	if err != nil {
		// Redirect to the login page if the cookie is missing
		c.Redirect(302, "/login")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for validation
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		// Redirect to the login page if token parsing fails
		c.Redirect(302, "/login")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Redirect(302, "/login")
			return
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		// Check if the user exists in the database
		if user.ID == 0 {
			c.Redirect(302, "/login")
			return
		}

		// Set user information in the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	} else {
		c.Redirect(302, "/login")
		return
	}
}
