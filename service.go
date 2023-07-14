package main

import (
	"errors"
)

type AuthServicer interface {
	SignUp(req *SingUpReq) error
	SignIn(req *SignInReq) (string, error)
}

type AuthService struct {
	store Store
}

func (as *AuthService) SignUp(req *SingUpReq) error {
	err := VerifyRegisterRequest(req)
	if err != nil {
		return err
	}
	req.Password, err = EncryptPassword(req.Password)
	if err != nil {
		return err
	}
	user := CreateUser(req)
	err = as.store.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) SignIn(req *SignInReq) (string, error) {
	mongouser, err := as.store.RetrieveUser(req.Email)
	if err != nil {
		return "", err
	}
	err = VerifyPassword(mongouser.Password, req.Password)
	if err != nil {
		return "", errors.New("incorrect password")
	}
	return "", nil
}
