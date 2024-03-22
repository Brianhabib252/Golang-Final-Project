package model

import (
	"time"
	//"golang-final-project/database"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// gorm model
type GormModel struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// User model
type User struct {
	GormModel
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required"`
}

// Photo model
type Photo struct {
	GormModel
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" validate:"required"`
	UserID   uint   `json:"user_id" gorm:"foreignKey:UserRefer"`
	User     User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// SocialMedia model
type SocialMedia struct {
	GormModel
	Name           string `json:"name" validate:"required"`
	SocialMediaURL string `json:"social_media_url" validate:"required"`
	UserID         uint   `json:"user_id" gorm:"foreignKey:UserRefer"`
	User           User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Comment model
type Comment struct {
	GormModel
	UserID  uint   `json:"user_id" gorm:"foreignKey:UserRefer"`
	User    User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PhotoID uint   `json:"photo_id" gorm:"foreignKey:PhotoRefer"`
	Photo   Photo  `json:"photo" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Message string `json:"message" validate:"required"`
}

// ValidateStruct validates the user struct according to the struct tags
func ValidateUser(user User) error {
	if err := validate.Struct(user); err != nil {
		return err
	}
	return nil
}

// ValidateStruct validates the photo struct according to the struct tags
func ValidatePhoto(photo Photo) error {
	if err := validate.Struct(photo); err != nil {
		return err
	}
	return nil
}

// ValidateStruct validates the social media struct according to the struct tags
func ValidateSosmed(sosmed SocialMedia) error {
	if err := validate.Struct(sosmed); err != nil {
		return err
	}
	return nil
}

// ValidateStruct validates the comment struct according to the struct tags
func ValidateComment(comment Comment) error {
	if err := validate.Struct(comment); err != nil {
		return err
	}
	return nil
}
