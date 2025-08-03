package user

import (
	"github.com/gofiber/fiber/v2"
	"user-service/pkg/config"
	"user-service/pkg/handler/utils"
)

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
	// 1. Cookie’den Refresh Token al
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh token yok"})
	}

	// 2. Refresh Token'ı doğrula
	claims, err := utils.ParseAccessToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh token geçersiz"})
	}

	// 3. Claims'ten email, user_id ve uuid al
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "geçersiz email"})
	}

	userID := int(claims["user_id"].(float64))
	uuid, ok := claims["uuid"].(string)
	if !ok || uuid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "geçersiz uuid"})
	}

	// 4. Redis'ten refresh token doğrula
	storedToken, _ := config.Rdb.Get(c.Context(), email).Result()
	if storedToken != refreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh token mismatch"})
	}

	// 5. Yeni Access Token oluştur
	newAccessToken, err := utils.GenerateAccessToken(userID, email, uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token üretilemedi"})
	}

	return c.JSON(fiber.Map{"access_token": newAccessToken})
}
