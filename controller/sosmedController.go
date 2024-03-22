package controller

import (
	"golang-final-project/middleware"
	"golang-final-project/model"
	"golang-final-project/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSocialMedia creates a new social media entry
func CreateSocialMedia(c *gin.Context) {
	var socialMedia model.SocialMedia
	socialMedia.UserID = middleware.IDJWT
	// Fetch user details based on provided user ID
	var user model.User
	if err := repository.GetByID(socialMedia.UserID, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}
	socialMedia.User = user
	// bind json request to struct
	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate social media struct
	if err := model.ValidateSosmed(socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create social media entry
	if err := repository.Create(&socialMedia); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media entry"})
		return
	}
	c.JSON(http.StatusCreated, socialMediaResponse(socialMedia))
}

// GetSocialMedia retrieves all social media entries
func GetSocialMedia(c *gin.Context) {
	var socialMedia []model.SocialMedia
	if err := repository.GetPreloadAll("User", &model.SocialMedia{}, &socialMedia); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social media entries"})
		return
	}
	var socialMediaResponses []gin.H
	for _, sm := range socialMedia {
		socialMediaResponses = append(socialMediaResponses, socialMediaResponse(sm))
	}
	c.JSON(http.StatusOK, socialMediaResponses)
}

// GetPhotoByID retrieves a single photo by ID
func GetSocialMediaByID(c *gin.Context) {
	id := c.Param("sosmed_id")
	// Convert string ID to uint
	sosmedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	// get photo request by photo id
	var socialMedia model.SocialMedia
	if err := repository.GetPreloadByID("User", uint(sosmedID), &socialMedia); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social Media not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Social Media"})
		return
	}
	c.JSON(http.StatusOK, socialMediaResponse(socialMedia))
}

// UpdateSocialMedia updates a social media entry by ID
func UpdateSocialMedia(c *gin.Context) {
	id := c.Param("sosmed_id")
	// Convert string ID to uint
	sosmedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	var socialMedia model.SocialMedia
	if err := repository.GetPreloadByID("User", uint(sosmedID), &socialMedia); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social media entry"})
		return
	}
	// bind reuest to struct
	var request struct {
		Name           string `json:"name"`
		SocialMediaUrl string `json:"social_media_url"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	socialMedia.Name = request.Name
	socialMedia.SocialMediaURL = request.SocialMediaUrl
	// validate the update social media struct
	if err := model.ValidateSosmed(socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update database
	if err := repository.Update(&socialMedia); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}
	c.JSON(http.StatusOK, socialMediaResponse(socialMedia))
}

// DeleteSocialMedia deletes a social media entry by ID
func DeleteSocialMedia(c *gin.Context) {
	id := c.Param("sosmed_id")
	// Convert string ID to uint
	sosmedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	var socialMedia model.SocialMedia
	if err := repository.GetByID(uint(sosmedID), &socialMedia); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social media entry"})
		return
	}
	if err := repository.Delete(&socialMedia); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Social media entry deleted successfully"})
}

// Helper function to create response with user details
func socialMediaResponse(socialMedia model.SocialMedia) gin.H {
	return gin.H{
		"social_media_id":  socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user": gin.H{
			"username": socialMedia.User.Username,
			"email":    socialMedia.User.Email,
		},
	}
}
