package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password,omitempty" binding:"required,min=6"`
	Age       int       `gorm:"type:int" json:"age" binding:"omitempty,gte=0,lte=150"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age" binding:"omitempty,gte=0,lte=150"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Age      int    `json:"age" binding:"omitempty,gte=0,lte=150"`
}
