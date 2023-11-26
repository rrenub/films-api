package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

type User struct {
	gorm.Model
	Name      string `gorm:"unique; not null"`
	Password  string `gorm:"not null"`
	Favorites []Favourite
	Movie     []Movie
}

func (m *UserModel) Authenticate(name, password string) (int, error) {
	var user User

	// Extract user from BD if exists
	result := m.DB.Where("name = ?", name).First(&user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether if password submitted hash is equal to saved password hash in DB
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, nil
		}
	}

	return int(user.ID), nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var user User

	r := m.DB.
		Where("`id` = ?", id).
		Limit(1).
		Find(&user)

	if err := r.Error; err != nil {
		return false, err
	}

	exists := r.RowsAffected > 0
	return exists, nil
}

func (m *UserModel) Insert(name, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := &User{
		Name:     name,
		Password: string(hashedPassword),
	}

	result := m.DB.Create(user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicatedEntry
		} else {
			return err
		}
	}

	return nil
}
