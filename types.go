package main

type SignupReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}
