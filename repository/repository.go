package repository

import (
	"golang-final-project/database"
	"golang-final-project/model"
	"time"

	"gorm.io/gorm"
)

// Create creates a new record in the database
func Create(data interface{}) error {
	db := database.GetDB()
	return db.Create(data).Error
}

// Get data by seraching user id in the tabble from the database
func GetByUserID(userID uint, model interface{}, records interface{}) error {
	db := database.GetDB()
	if err := db.Find(records, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

// GetByUserIDAndDate retrieves user nutritions based on user ID and date from the database
func GetByUserIDAndDate(date time.Time, userID uint, model interface{}, records interface{}) error {
	db := database.GetDB()
	if err := db.Where("user_id = ? AND DATE(created_at) = ?", userID, date.Format("2006-01-02")).Find(records).Error; err != nil {
		return err
	}
	return nil
}

// Get Preload All retrieves all records of a model from the database
func GetPreloadAll(preload string, model interface{}, records interface{}) error {
	db := database.GetDB()
	if err := db.Preload(preload).Find(records).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a single record of a model by ID from the database
func GetByID(id uint, model interface{}) error {
	db := database.GetDB()
	if err := db.First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// GetPreloadByID retrieves a single record of a model by ID from the database
func GetPreloadByID(load string, id uint, model interface{}) error {
	db := database.GetDB()
	if err := db.Preload(load).First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Get preload data retrieves a single record of a model by ID from the database
func GetPreload(load string, model interface{}) error {
	db := database.GetDB()
	if err := db.Preload(load).Find(model).Error; err != nil {
		return err
	}
	return nil
}

// Get two preload data retrieves a single record of a model by ID from the database
func GetTwoPreloadByID(load1 string, load2 string, id uint, model interface{}) error {
	db := database.GetDB()
	if err := db.Preload(load1).Preload(load2).First(model, id).Error; err != nil {
		return err
	}
	return nil
}

// Get Nutrition data retrieves a single record of a model by food name from the database
func GetByfoodName(foodName string, model interface{}) error {
	db := database.GetDB()
	if err := db.First(&model, "name_food = ?", foodName).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a single record of a model by ID from the database
func GetByEmail(email string, model interface{}) error {
	db := database.GetDB()
	if err := db.First(&model, "email = ?", email).Error; err != nil {
		return err
	}
	return nil
}

// Update updates a record in the database
func Update(updatedData interface{}) error {
	db := database.GetDB()
	if err := db.Save(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by ID
func Delete(model interface{}) error {
	db := database.GetDB()
	if err := db.Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func IsEmailUnique(email string) bool {
	db := database.GetDB()
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		return false // Email already exists
	} else if err != gorm.ErrRecordNotFound {
		return false // Database error
	}
	return true // Email is unique
}

func IsUsernameUnique(username string) bool {
	db := database.GetDB()
	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err == nil {
		return false // Username already exists
	} else if err != gorm.ErrRecordNotFound {
		return false // Database error
	}
	return true // Username is unique
}
