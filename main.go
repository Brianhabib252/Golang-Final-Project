package main

import (
	//"fmt"
	"golang-final-project/controller"
	"golang-final-project/database"
	"golang-final-project/middleware"

	//"golang-final-project/model"
	"github.com/gin-gonic/gin"
)

func main() {
	database.StartdB()
	database.AutoMigrate()

	r := gin.Default()

	// User Routes
	r.POST("/signup", controller.SignUp)
	r.POST("/signin", controller.SignIn)
	r.PUT("/user/:id", middleware.Authentication, middleware.Authorization(), controller.UpdateUsername)
	r.DELETE("/user/:id", middleware.Authentication, middleware.Authorization(), controller.DeleteUser)
	r.POST("/signout", middleware.Authentication, controller.SignOut)

	// Photo Routes
	r.POST("/photo", middleware.Authentication, controller.CreatePhoto)
	r.GET("/photo", middleware.Authentication, controller.GetPhotos)
	r.GET("/photo/:photo_id", middleware.Authentication, controller.GetPhotoByID)
	r.PUT("/photo/:id/:photo_id", middleware.Authentication, middleware.Authorization(), controller.UpdatePhoto)
	r.DELETE("/photo/:id/:photo_id", middleware.Authentication, middleware.Authorization(), controller.DeletePhoto)

	// Sosmed Routes
	r.POST("/sosmed", middleware.Authentication, controller.CreateSocialMedia)
	r.GET("/sosmed", middleware.Authentication, controller.GetSocialMedia)
	r.GET("/sosmed/:sosmed_id", middleware.Authentication, controller.GetSocialMediaByID)
	r.PUT("/sosmed/:id/:sosmed_id", middleware.Authentication, middleware.Authorization(), controller.UpdateSocialMedia)
	r.DELETE("/sosmed/:id/:sosmed_id", middleware.Authentication, middleware.Authorization(), controller.DeleteSocialMedia)

	// Comment Routes
	r.POST("/comment", middleware.Authentication, controller.CreateComment)
	r.GET("/comment/:comment_id", middleware.Authentication, controller.GetCommentByID)
	r.PUT("/comment/:id/:comment_id", middleware.Authentication, middleware.Authorization(), controller.UpdateComment)
	r.DELETE("/comment/:id/:comment_id", middleware.Authentication, middleware.Authorization(), controller.DeleteComment)

	// Run the server
	r.Run(":8080")
}
