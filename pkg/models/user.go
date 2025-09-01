package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// User struct now includes the Password field for authentication
type User struct {
	gorm.Model
	Name         string  `json:"name"`
	Email        string  `json:"email" gorm:"unique"`
	Password     string  `json:"password,omitempty"` // omitempty prevents it from being sent in responses if empty
	MembershipID string  `json:"membership_id" gorm:"unique"`
	Role         string  `json:"role"` // "staff" or "student"
	Fines        float64 `json:"fines"`
}

// BeforeSave is a GORM hook that automatically hashes the password before saving a user.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return
}

// CreateUser creates a new user in the database.
func (u *User) CreateUser() (*User, error) {
	result := db.Create(&u)
	return u, result.Error
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
}

// GetUserById retrieves a user by their ID.
func GetUserById(Id int64) (*User, *gorm.DB) {
	var getUser User
	result := db.Where("id = ?", Id).First(&getUser)
	return &getUser, result
}

// DeleteUser deletes a user by their ID.
func DeleteUser(Id int64) User {
	var user User
	db.Where("id = ?", Id).Delete(&user)
	return user
}
