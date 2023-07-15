package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New(fiber.Config{
		PassLocalsToViews: true,
		Views:             createEngine(),
	})

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Print(err)
	}

	store := &MongoStore{db: client.Database("users"), coll: "users"}
	service := &AuthService{store: store}
	as := &apiServer{s: service}

	app.Get("/", as.home)
	app.Get("/signup", loggingMiddleware(as.handleSignUp))
	app.Post("/signup", loggingMiddleware(as.handleSignUp))
	app.Post("/signin", loggingMiddleware(as.handleSignIn))

	log.Fatal(app.Listen(":3000"))
}

type apiServer struct {
	s AuthServicer
}

func (as *apiServer) home(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

func (as *apiServer) handleSignUp(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Render("signupform", nil)
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	req := &SignupReq{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := as.s.SignUp(req); err != nil {
		return err
	}

	return c.Render("index", nil)
}

func (as *apiServer) handleSignIn(c *fiber.Ctx) error {

	req := &SignInReq{}
	if err := json.Unmarshal(c.Body(), req); err != nil {
		return err
	}

	username, err := as.s.SignIn(req)
	if err != nil {
		return err
	}

	token, err := createJWT(username)
	if err != nil {
		return err
	}

	c.Status(200)
	return c.JSON(token)
}

func createJWT(username string) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(15 * time.Minute).Unix(),
		"username":  username,
	}

	secret := "secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func loggingMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("req on path", c.Path(), "method", c.Method())
		return next(c)
	}
}
