package models

import "gorm.io/gorm"

// User struct defines the user model
type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
