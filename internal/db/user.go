package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Email     string     `json:"email" gorm:"unique"`
	Password  string     `json:"-"`
	IsAdmin  bool       `json:"isAdmin" gorm:"default:false"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
	user := User{
		Email:    "test1@mail.com",
		Password: "123123",
		IsAdmin: true,
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return errors.New("error creating password")
	}
	user.Password = string(password)

	if err := DBConn.Create(&user).Error; err != nil {
		return errors.New("error creating user")
	}
	return nil
}

func (u *User) LoginAsAdmin(email, password string) (*User, error) {
	// Find
	query := "email = ? AND is_admin = ?"
	if err := DBConn.Where(query, email, true).First(&u).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return u, nil
}
