package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"text/template"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	home := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		data := map[string]string{
			"wow": "mega",
		}
		tmpl.Execute(w, data)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Print(err)
	}

	store := &MongoStore{db: client.Database("users"), coll: "users"}
	service := &AuthService{store: store}
	as := &apiServer{s: service}

	http.HandleFunc("/", home)
	http.HandleFunc("/signup", makeHTTPHandleFunc(as.handleSignUp))
	http.HandleFunc("/signin", makeHTTPHandleFunc(as.handleSignIn))

	log.Fatal(http.ListenAndServe(":3000", nil))
}

type apiServer struct {
	s AuthServicer
}

func (as *apiServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return errors.New("method not allowed")
	}

	req := &SingUpReq{}
	json.NewDecoder(r.Body).Decode(req)

	err := as.s.SignUp(req)
	if err != nil {
		return err
	}

	return WriteResponse(w, 200, "succesfully register user")
}

func (as *apiServer) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return errors.New("method not allowed")
	}

	req := &SignInReq{}
	json.NewDecoder(r.Body).Decode(&req)
	username, err := as.s.SignIn(req)
	if err != nil {
		return err
	}

	token, err := createJWT(username)
	if err != nil {
		return err
	}
	return WriteResponse(w, 200, token)
}

func WriteResponse(w http.ResponseWriter, code int, data ...any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(&data)
}

func createJWT(username string) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"username":  username,
	}

	secret := "secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteResponse(w, http.StatusBadRequest, "err: "+err.Error())
		}
	}
}
