package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"golang-final-project/middleware"
	"golang-final-project/model"
	"golang-final-project/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Nutrition struct {
	Energy  float64
	Protein float64
	Carbo   float64
	Fat     float64
	Fiber   float64
}

// CreateUserNutrition creates a new user nutrition entry
func CreateUserNutrition(c *gin.Context) {
	var userNutrition model.UserNutrition
	userNutrition.UserID = middleware.IDJWT
	// bind request to struct
	if err := c.ShouldBindJSON(&userNutrition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get social media from database
	var socialMedia model.SocialMedia
	if err := repository.GetPreloadByID("User", userNutrition.SocialMediaID, &socialMedia); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}
	userNutrition.SocialMedia = socialMedia
	// validate nutrition user struct
	if err := model.ValidateUserNutrition(userNutrition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// create data to database
	if err := repository.Create(&userNutrition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user nutrition"})
		return
	}

	c.JSON(http.StatusCreated, userNutritionResponse(userNutrition))
}

// GetUserNutritions retrieves all user nutritions
func GetUserNutritions(c *gin.Context) {
	// get nutrition data from database
	var userNutritions []model.UserNutrition
	userID := middleware.IDJWT
	if err := repository.GetByUserID(userID, &model.UserNutrition{}, &userNutritions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user nutritions"})
		return
	}
	var userNutritionResponses []gin.H
	fmt.Println(userNutritions)
	for _, userNutrition := range userNutritions {
		userNutritionResponses = append(userNutritionResponses, userNutritionResponse(userNutrition))
	}
	// Extract food names from user nutritions
	var foodNames []string
	for _, userNutrition := range userNutritions {
		foodNames = append(foodNames, userNutrition.Food)
	}
	fmt.Println(foodNames)
	c.JSON(http.StatusOK, userNutritionResponses)
}

// GetUserNutritionsByDate retrieves all user nutritions based on a day date
func GetUserNutritionsByDate(c *gin.Context) {
	id := c.Param("sosmed_id")
	// Convert string ID to uint
	sosmedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	// Parse date from request query parameters
	dateString := c.Query("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	var userNutritions []model.UserNutrition
	userID := middleware.IDJWT
	if err := repository.GetByUserIDAndDate(date, userID, &model.UserNutrition{}, &userNutritions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user nutritions"})
		return
	}
	// Extract food names from user nutritions
	var foodNames []string
	var foodMass []int
	for _, userNutrition := range userNutritions {
		foodNames = append(foodNames, userNutrition.Food)
		foodMass = append(foodMass, userNutrition.FoodMass)
	}
	var nutritionData model.NutritionData
	var nutritionUser Nutrition
	for i, food := range foodNames {
		if err := repository.GetByfoodName(food, &nutritionData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email or Password"})
			return
		} else {
			nutritionUser.Energy += nutritionData.Energy * float64(foodMass[i]) / 100
			nutritionUser.Carbo += nutritionData.Carbo * float64(foodMass[i]) / 100
			nutritionUser.Protein += nutritionData.Protein * float64(foodMass[i]) / 100
			nutritionUser.Fat += nutritionData.Fat * float64(foodMass[i]) / 100
			nutritionUser.Fiber += nutritionData.Fiber * float64(foodMass[i]) / 100
		}
	}
	nutritionNeed, err := NutritionCaloriesCalculation(uint(sosmedID))
	if err != nil {
		// Handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate nutrition"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"energy":       nutritionUser.Energy,
		"energy_need":  nutritionNeed.Energy,
		"protein":      nutritionUser.Protein,
		"protein_need": nutritionNeed.Protein,
		"carbo":        nutritionUser.Carbo,
		"carbo_need":   nutritionNeed.Carbo,
		"fat":          nutritionUser.Fat,
		"fat_need":     nutritionNeed.Fat,
		"fiber":        nutritionUser.Fiber,
		"fiber_need":   nutritionNeed.Fiber,
	})
}

// DeleteUserNutrition deletes a user nutrition entry by ID
func DeleteUserNutrition(c *gin.Context) {
	id := c.Param("nutrition_id")
	nutritionID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	var userNutrition model.UserNutrition
	if err := repository.GetByID(uint(nutritionID), &userNutrition); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User nutrition not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user nutrition"})
		return
	}

	if err := repository.Delete(&userNutrition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user nutrition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User nutrition deleted successfully"})
}

func userNutritionResponse(userNutrition model.UserNutrition) gin.H {
	return gin.H{
		"nutrition_id": userNutrition.ID,
		"user_id":      userNutrition.UserID,
		"food":         userNutrition.Food,
		"food_mass":    userNutrition.FoodMass,
		"date":         userNutrition.CreatedAt,
	}
}

func NutritionCaloriesCalculation(sosmedID uint) (Nutrition, error) {
	var nutrition Nutrition

	// Fetch user details from the database
	var socialMedia model.SocialMedia
	err := repository.GetPreloadByID("User", sosmedID, &socialMedia)
	if err != nil {
		return nutrition, err
	}

	age := float64(socialMedia.User.Age)
	height := float64(socialMedia.Height)
	weight := float64(socialMedia.Weight)
	gender := socialMedia.Gender

	// Daily Calories Needed
	var energy float64
	if gender == "MALE" {
		energy = math.Round(66 + (13.7 * weight) + (5 * height) - (6.8 * age))
	} else {
		energy = math.Round(655 + (9.6 * weight) + (1.8 * height) - (4.7 * age))
	}

	// Nutrient Calculation (Gram)
	protein := math.Round((15.0 / 100.0 * energy) / 4.0)
	carbo := math.Round((60.0 / 100.0 * energy) / 4.0)
	fat := math.Round((15.0 / 100.0 * energy) / 9.0)
	fiber := age + 10

	nutrition = Nutrition{
		Energy:  energy,
		Protein: protein,
		Carbo:   carbo,
		Fat:     fat,
		Fiber:   fiber,
	}

	return nutrition, nil
}
