package main

import (
	"golang-final-project/controller"
	"golang-final-project/database"
	"golang-final-project/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.StartdB()
	database.AutoMigrate()

	routes := gin.Default()

	// User Routes
	routes.POST("/signup", controller.SignUp)
	routes.POST("/signin", controller.SignIn)
	routes.PUT("/user/:id", middleware.Authentication, middleware.Authorization(), controller.UpdateUsername)
	routes.DELETE("/user/:id", middleware.Authentication, middleware.Authorization(), controller.DeleteUser)
	routes.POST("/signout", middleware.Authentication, controller.SignOut)

	// Group routes that require authentication
	authRoutes := routes.Group("/", middleware.Authentication)

	// Photo Routes
	authRoutes.POST("/photo", controller.CreatePhoto)
	authRoutes.GET("/photo", controller.GetPhotos)
	authRoutes.GET("/photo/:photo_id", controller.GetPhotoByID)
	authRoutes.PUT("/photo/:id/:photo_id", middleware.Authorization(), controller.UpdatePhoto)
	authRoutes.DELETE("/photo/:id/:photo_id", middleware.Authorization(), controller.DeletePhoto)

	// Sosmed Routes
	authRoutes.POST("/sosmed", controller.CreateSocialMedia)
	authRoutes.GET("/sosmed", controller.GetSocialMedia)
	authRoutes.GET("/sosmed/:sosmed_id", controller.GetSocialMediaByID)
	authRoutes.PUT("/sosmed/:id/:sosmed_id", middleware.Authorization(), controller.UpdateSocialMedia)
	authRoutes.DELETE("/sosmed/:id/:sosmed_id", middleware.Authorization(), controller.DeleteSocialMedia)

	// Comment Routes
	authRoutes.POST("/comment", controller.CreateComment)
	authRoutes.GET("/comment/:comment_id", controller.GetCommentByID)
	authRoutes.PUT("/comment/:id/:comment_id", middleware.Authorization(), controller.UpdateComment)
	authRoutes.DELETE("/comment/:id/:comment_id", middleware.Authorization(), controller.DeleteComment)

	// User Nutrition Routes
	authRoutes.POST("/nutrition", controller.CreateUserNutrition)
	authRoutes.GET("/nutrition", controller.GetUserNutritions)
	authRoutes.GET("/nutrition/:sosmed_id/queri", controller.GetUserNutritionsByDate)
	authRoutes.DELETE("/nutrition/:id/:nutrition_id", middleware.Authorization(), controller.DeleteUserNutrition)

	// Run the server
	routes.Run(":8080")
}
