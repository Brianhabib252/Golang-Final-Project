package controller

import (
	"fmt"
	"golang-final-project/middleware"
	"golang-final-project/model"
	"golang-final-project/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePhoto creates a new photo
func CreatePhoto(c *gin.Context) {
	var photo model.Photo
	photo.UserID = middleware.IDJWT
	// Fetch user details based on provided user ID
	var user model.User
	if err := repository.GetByID(photo.UserID, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}
	photo.User = user
	// get request body
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// validate photo struct
	if err := model.ValidatePhoto(photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// create photo request on database
	if err := repository.Create(&photo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo"})
		return
	}
	c.JSON(http.StatusCreated, photoResponse(photo))
}

// GetPhotos retrieves all photos
func GetPhotos(c *gin.Context) {
	var photos []model.Photo
	// get photo from database
	if err := repository.GetPreloadAll("User", &model.Photo{}, &photos); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photos"})
		return
	}
	// get user from database
	var photoResponses []gin.H
	for _, photo := range photos {
		photoResponses = append(photoResponses, photoResponse(photo))
	}
	c.JSON(http.StatusOK, photoResponses)
}

// GetPhotoByID retrieves a single photo by ID
func GetPhotoByID(c *gin.Context) {
	id := c.Param("photo_id")
	photoID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	// get photo request by photo id
	var photo model.Photo
	if err := repository.GetPreloadByID("User", uint(photoID), &photo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photo"})
		return
	}
	c.JSON(http.StatusOK, photoResponse(photo))
}

// UpdatePhoto updates a photo by ID
func UpdatePhoto(c *gin.Context) {
	id := c.Param("photo_id")
	photoID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	// get photo request by photo id
	var photo model.Photo
	photo.UserID = middleware.IDJWT
	if err := repository.GetByID(uint(photoID), &photo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}
	// get photo request from database
	if err := repository.GetPreloadByID("User", uint(photoID), &photo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photo"})
		return
	}
	// make request struct for update
	var request struct {
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photo_url"`
	}
	// bind update photo to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	photo.Title = request.Title
	photo.Caption = request.Caption
	// validate photo struct
	if err := model.ValidatePhoto(photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// update database
	if err := repository.Update(&photo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error to update photo"})
		return
	}
	c.JSON(http.StatusOK, photoResponse(photo))
}

// DeletePhoto deletes a photo by ID
func DeletePhoto(c *gin.Context) {
	id := c.Param("photo_id")
	photoID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	// check photo id in database
	var photo model.Photo
	if err := repository.GetByID(uint(photoID), &photo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photo"})
		return
	}
	// delete photo
	if err := repository.Delete(&photo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error to delete photo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

// Helper function to create response with user details
func photoResponse(photo model.Photo) gin.H {
	return gin.H{
		"photo_id":   photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"created_at": photo.CreatedAt,
		"updated_at": photo.UpdatedAt,
		"user": gin.H{
			"username": photo.User.Username,
			"email":    photo.User.Email,
		},
	}
}
