package Auth

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"simpleTaskManager/MongoDB"
	"strings"
	"time"
)

var (
	Store    *session.Store
	AUTH_KEY string = "authenticated"
	USER_ID  string = "USER_ID"
)

func InitStore(s *session.Store) {
	Store = s
}

func NewMiddleware() fiber.Handler {
	return AuthMiddleware
}

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := Store.Get(c)

	pathSegments := strings.Split(c.Path(), "/")
	if len(pathSegments) > 1 && pathSegments[1] == "auth" {
		return c.Next()
	}

	if err != nil || sess.Get(AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	return c.Next()
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type RegisterLoginUser struct {
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

func CreateUser(user *RegisterLoginUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var checkUser User
	err := MongoDB.AuthCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&checkUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if user.Email == checkUser.Email {
		return fmt.Errorf("email already registered")
	}

	_, dbErr := MongoDB.AuthCollection.InsertOne(ctx, user)
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func GetUser(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := MongoDB.AuthCollection.FindOne(ctx, bson.M{"_id": objID})
	var user User
	err = result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckEmail(email string, user *User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := MongoDB.AuthCollection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil || email != user.Email {
		return false
	}
	return true
}

func Register(c *fiber.Ctx) error {

	var data RegisterLoginUser

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + err.Error(),
		})
	}

	password, bcErr := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if bcErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + bcErr.Error(),
		})
	}

	user := RegisterLoginUser{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(password),
	}
	err = CreateUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "registered.",
	})
}

func Login(c *fiber.Ctx) error {
	var data RegisterLoginUser

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + err.Error(),
		})
	}

	var user User
	if !CheckEmail(data.Email, &user) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}
	sess, sessErr := Store.Get(c)
	if sessErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + sessErr.Error(),
		})
	}

	sess.Set(AUTH_KEY, true)
	sess.Set(USER_ID, (user.ID).Hex())

	sessErr = sess.Save()
	if sessErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + sessErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged in",
	})
}

func Logout(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "logged out (no session)",
		})
	}

	err = sess.Destroy()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out",
	})
}

func HealthCheck(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	auth := sess.Get(AUTH_KEY)
	if auth != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "authenticated",
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}
}
