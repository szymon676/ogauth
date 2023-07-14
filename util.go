package main

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req *SingUpReq) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: NormalizeUsername(req.Username),
		Email:    NormalizeEmail(req.Email),
		Password: req.Password,
	}
}

func VerifyRegisterRequest(req *SingUpReq) error {
	if len(req.Username) < 2 {
		return errors.New("username must be longer than 2 characters")
	}
	if len(req.Email) < 4 {
		return errors.New("email must be longer than 4 characters")
	}
	if len(req.Password) < 4 {
		return errors.New("password must be longer than 4 characters")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
