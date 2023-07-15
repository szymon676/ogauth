package main

import (
	"context"
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
	as := &apiServer{s: service, signed: false}

	app.Get("/", as.home)
	app.Get("/signup", loggingMiddleware(as.handleSignUpPage))
	app.Get("/signin", loggingMiddleware(as.handleSignInPage))
	app.Post("/signup", loggingMiddleware(as.handleSignUp))
	app.Post("/signin", loggingMiddleware(as.handleSignIn))
	app.Post("/signout", loggingMiddleware(as.handleSignOut))

	log.Fatal(app.Listen(":3000"))
}

type apiServer struct {
	s      AuthServicer
	signed bool
	token  string
}

func (as *apiServer) home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{"logged": as.signed})
}

func (as *apiServer) handleSignUpPage(c *fiber.Ctx) error {
	return c.Render("signupform", nil)
}

func (as *apiServer) handleSignInPage(c *fiber.Ctx) error {
	return c.Render("signinform", nil)
}

func (as *apiServer) handleSignUp(c *fiber.Ctx) error {
	var req SignupReq
	c.BodyParser(&req)

	if err := as.s.SignUp(&req); err != nil {
		return err
	}

	return c.Redirect("/signin")
}

func (as *apiServer) handleSignIn(c *fiber.Ctx) error {
	var req SignInReq
	c.BodyParser(&req)

	username, err := as.s.SignIn(&req)
	if err != nil {
		return err
	}

	token, err := createJWT(username)
	if err != nil {
		return err
	}

	as.token = token
	as.signed = true

	return c.SendString("signed in")
}

func (as *apiServer) handleSignOut(c *fiber.Ctx) error {
	as.signed = false
	as.token = ""
	return c.SendString("signed out")
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
