package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"user-service/pkg/config"
	"user-service/pkg/handler/user/model"
)

// JWT Secret .env'den al
func GetJWTSecret() []byte {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	return []byte(os.Getenv("JWT_SECRET"))
}

// Access Token oluştur
func GenerateAccessToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

// Refresh Token oluştur
func GenerateRefreshToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "refresh",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

// Token çözümle
func ParseAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	return token.Claims.(jwt.MapClaims), nil
}

// Refresh Token Cookie Set
func SetRefreshToken(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Prod'da true
		SameSite: "Lax",
	})
}

// Refresh Token Sil
func ClearSetRefreshToken(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
}
func GetUserByID(userID int) (*model.User, error) {
	var user model.User
	err := config.DB.QueryRow("SELECT id, name, email, age, password, role FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
