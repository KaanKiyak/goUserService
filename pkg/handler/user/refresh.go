package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"user-service/pkg/config"
	"user-service/pkg/handler/utils"
)

/* func RefreshHandle(c *fiber.Ctx) error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	cookie := c.Cookies("refresh_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	var email string
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	claims := token.Claims.(jwt.MapClaims)
	email = claims["email"].(string)

	storedToken, _ := config.Rdb.Get(c.Context(), email).Result()
	if storedToken != cookie {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	newAccessToken, err := utils.GenerateAccessToken(fmt.Sprintf("%d", int(claims["user_id"].(float64))), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "token oluşturulamadı",
		})
	}

	return c.JSON(fiber.Map{"access_token": newAccessToken})
}*/

// RefreshHandle access token yeniler
// @Summary Access Token Yenileme
// @Description Geçerli refresh token kullanılarak yeni bir access token üretilir
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Yeni access token döner"
// @Failure 401 {object} ErrorResponse "Refresh token yok veya geçersiz"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /refresh [post]
func RefreshHandle(c *fiber.Ctx) error {
	cookie := c.Cookies("refresh_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no refresh token"})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	userID := int(claims["user_id"].(float64))

	storedToken, _ := config.Rdb.Get(c.Context(), email).Result()
	if storedToken != cookie {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh token mismatch"})
	}

	newAccessToken, err := utils.GenerateAccessToken(userID, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate new token"})
	}

	return c.JSON(fiber.Map{"access_token": newAccessToken})
}
