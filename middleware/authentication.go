package middleware

import (
	"fmt"
	"golang-final-project/model"
	"golang-final-project/repository"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var IDJWT uint

func Authentication(c *gin.Context) {
	// get database variable from environment
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("cannot get enc file :")
	}
	secretKey := os.Getenv("SECRET_KEY")
	// get request cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Println("Error retrieving token from cookie")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token not found",
		})
		c.Abort()
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println("Error Validating JWT :", err)
		c.Abort()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check token exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// get User Id from jwt
		userIDJWT, _ := claims["sub"].(float64)
		IDJWT = uint(userIDJWT)
		// check user
		var user model.User
		err = repository.GetByID(IDJWT, &user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}
		// Attach to req and continue
		c.Set("user", user)
		c.Next()
	} else {
		fmt.Println(err)
	}
}
