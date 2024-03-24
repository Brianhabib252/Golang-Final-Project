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

// CreateComment creates a new comment
func CreateComment(c *gin.Context) {
	var comment model.Comment
	comment.UserID = middleware.IDJWT
	// Fetch user details based on provided user ID
	var user model.User
	if err := repository.GetByID(comment.UserID, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comment details"})
		return
	}
	comment.User = user
	// bind request to struct
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Fetch user details based on provided user ID
	var photo model.Photo
	if err := repository.GetPreloadByID("User", comment.PhotoID, &photo); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}
	comment.Photo = photo
	// Validate comment struct
	if err := model.ValidateComment(comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create comment entry
	if err := repository.Create(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, commentResponse(comment))
}

// GetCommentByID retrieves a single comment by ID
func GetCommentByID(c *gin.Context) {
	id := c.Param("comment_id")
	// Convert string ID to uint
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
		return
	}
	var comment model.Comment
	if err := repository.GetTwoPreloadByID("User", "Photo", uint(commentID), &comment); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}
	c.JSON(http.StatusOK, commentResponse(comment))
}

// UpdateComment updates a comment by ID
func UpdateComment(c *gin.Context) {
	id := c.Param("comment_id")
	// Convert string ID to uint
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
		return
	}
	var comment model.Comment
	comment.UserID = middleware.IDJWT
	// get comment from the database
	if err := repository.GetTwoPreloadByID("Photo", "User", uint(commentID), &comment); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}
	var request struct {
		PhotoID int    `json:"photo_id"`
		Message string `json:"message"`
	}
	// bind request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment.Message = request.Message
	comment.Photo.User = comment.User
	fmt.Println(comment)
	// update commant to database
	if err := model.ValidateComment(comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update database
	if err := repository.Update(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to update comment"})
		return
	}
	c.JSON(http.StatusOK, commentResponse(comment))
}

// DeleteComment deletes a comment by ID
func DeleteComment(c *gin.Context) {
	id := c.Param("comment_id")
	// Convert string ID to uint
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
		return
	}
	var comment model.Comment
	if err := repository.GetByID(uint(commentID), &comment); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}
	// delete comment
	if err := repository.Delete(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// Helper function to create response with comment details
func commentResponse(comment model.Comment) gin.H {
	return gin.H{
		"comment_id": comment.ID,
		"message":    comment.Message,
		"user": gin.H{
			"username": comment.User.Username,
			"email":    comment.User.Email,
		},
		"photo": gin.H{
			"title":     comment.Photo.Title,
			"caption":   comment.Photo.Caption,
			"photo_url": comment.Photo.PhotoURL,
		},
	}
}
