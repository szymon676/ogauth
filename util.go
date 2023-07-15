package main

import (
	"errors"
	"strings"

	"github.com/gofiber/template/django/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createEngine() *django.Engine {
	engine := django.New("./views", ".html")
	engine.Reload(true)
	return engine
}

func CreateUser(req *SignupReq) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: NormalizeUsername(req.Username),
		Email:    NormalizeEmail(req.Email),
		Password: req.Password,
	}
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func VerifyRegisterRequest(req *SignupReq) error {
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

func VerifyPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
