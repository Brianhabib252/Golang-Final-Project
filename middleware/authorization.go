package middleware

import (
	"fmt"
	"golang-final-project/model"
	"golang-final-project/repository"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get database variable from environment
		var err error
		err = godotenv.Load()
		if err != nil {
			fmt.Println("cannot get enc file :")
		}
		secretKey := os.Getenv("SECRET_KEY")
		// Get the JWT token from the cookie
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found in cookie"})
			c.Abort()
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method of the token
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Replace "YOUR_SECRET_KEY" with your actual JWT secret key
			return []byte(secretKey), nil
		})
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}
		// Extract user ID from the JWT claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}
		userID, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in token"})
			c.Abort()
			return
		}

		// Get resource ID from the request (for example, from path parameter)
		id := c.Param("id")
		// Convert string ID to uint
		resourceID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
			return
		}

		// Check if the user is the creator of the resource
		var resource model.User
		if err := repository.GetByID(uint(resourceID), &resource); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			c.Abort()
			return
		}
		if uint(userID) != resource.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to perform this action"})
			c.Abort()
			return
		}
		// Proceed to the next middleware or handler
		c.Next()
	}
}
