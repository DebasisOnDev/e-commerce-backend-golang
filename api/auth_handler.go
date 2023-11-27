package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type LoginUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterUserParams struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confPassword"`
}

const bcryptCost = 10

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CreateUserFromRegisterParams(registerParams RegisterUserParams) (*types.User, error) {
	hashedPassword, err := HashPassword(registerParams.Password)
	if err != nil {
		return nil, err
	}
	return &types.User{
		FirstName:         registerParams.FirstName,
		LastName:          registerParams.LastName,
		Email:             registerParams.Email,
		EncryptedPassword: hashedPassword,
		IsAdmin:           false,
	}, nil
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"-"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (a *AuthHandler) HandleRegisterUser(c *fiber.Ctx) error {
	var params RegisterUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if params.Password != params.ConfirmPassword {
		fmt.Printf("The password is %s", params.Password)
		fmt.Printf("The confirmed password is %s", params.ConfirmPassword)
		return c.JSON(map[string]string{"error": "the password and confirm password must be same"})
	}
	_, err := a.userStore.GetUserByEmail(c.Context(), params.Email)
	if err == nil {
		return c.JSON(map[string]string{"error": "user already exists"})
	}
	user, err := CreateUserFromRegisterParams(params)
	if err != nil {
		return err
	}
	newUser, err := a.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return invalidCredentials(c)
	}
	return fmt.Errorf("user is created with %s", newUser.ID)
}

func (a *AuthHandler) HandleLogInUser(c *fiber.Ctx) error {
	var params LoginUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := a.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(fiber.Map{"error": "document not found"})
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return c.JSON(fiber.Map{"error": "invalid password"})
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	c.Set("X-Api-Token", resp.Token)
	return c.JSON(resp)
}

func (a *AuthHandler) HandleLogOutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "X-Api-Token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 12).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
