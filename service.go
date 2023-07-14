package main

type AuthServicer interface {
	SignUp(req *SingUpReq) error
	Login(req *SignInReq) error
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
