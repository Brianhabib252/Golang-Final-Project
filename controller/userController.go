package controller

import (
	"fmt"
	"golang-final-project/model"
	"golang-final-project/repository"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get request body
	var user model.User
	// Bind request body to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request body"})
		return
	}
	//Validate user struct
	if err := model.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if the email is already registered
	if !repository.IsEmailUnique(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	// Check if the username is already registered
	if !repository.IsUsernameUnique(user.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	// Set the hashed password
	user.Password = string(hash)
	// Create user
	if err := repository.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	// Fetch user data from the database
	var newUser model.User
	if err := repository.GetByID(user.ID, &newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get data from database"})
		return
	}
	// respon
	c.JSON(http.StatusOK, gin.H{
		"id":       newUser.ID,
		"username": newUser.Username,
		"email":    newUser.Email,
		"age":      newUser.Age,
	})
}

func SignIn(c *gin.Context) {
	// get database variable from environment
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("cannot get enc file :")
	}
	secretKey := os.Getenv("SECRET_KEY")
	// get request body
	var body struct {
		Email    string
		Password string
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to req body",
		})
		return
	}
	//look up the request user
	var user model.User
	if err := repository.GetByEmail(body.Email, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email or Password"})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}
	// compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}
	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("error :", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to Create Token",
		})
		return
	}
	// send token as response
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"user_id": user.ID,
	})
}

func UpdateUsername(c *gin.Context) {
	// preload the items from order tables by id
	id := c.Param("id")
	// Convert string ID to uint
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	var user model.User
	if err := repository.GetByID(uint(userID), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get data from database"})
		return
	}
	// make request struct for update
	var request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Check if the username is already registered
	if !repository.IsUsernameUnique(request.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	// Update order fields
	user.Username = request.Username
	//Validate user struct
	if err := model.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update the order to the database
	if err := repository.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update database"})
		return
	}
	// respon
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"age":      user.Age,
	})
}

func DeleteUser(c *gin.Context) {
	// get preload item from order table by id
	id := c.Param("id")
	// Convert string ID to uint
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	var user model.User
	if err := repository.GetByID(uint(userID), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	// Delete order from database
	if err := repository.Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}
	c.JSON(200, gin.H{"message": "Order deleted successfully"})
}

func SignOut(c *gin.Context) {
	// get request body
	var body struct {
		Email    string
		Password string
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to req body",
		})
		return
	}
	//look up the request user
	var user model.User
	if err := repository.GetByEmail(body.Email, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email or Password"})
		return
	}
	// compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Password",
		})
		return
	}
	// black list token
	var tokensBlacklist = map[string]bool{}
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token not found",
		})
		return
	}
	tokensBlacklist[tokenString] = true // Add token to blacklist
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
