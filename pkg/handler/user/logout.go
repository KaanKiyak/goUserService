package user

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
	"user-service/pkg/config"
	"user-service/pkg/handler/utils"
)

// LogoutHandler kullanıcının oturumunu kapatır
// @Summary Kullanıcı Çıkışı (Logout)
// @Description Refresh token ve access token doğrulanır, Redis'ten refresh token silinir, cookie temizlenir
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout başarılı"
// @Failure 401 {object} ErrorResponse "Token bulunamadı veya geçersiz"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /logout [post]
func LogoutHandler(c *fiber.Ctx) error {
	// 1. Cookie'den Refresh Token al
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token bulunamadı",
		})
	}

	// 2. Authorization Header'dan Access Token al (isteğe bağlı)
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "geçersiz token",
		})
	}

	// 3. Email bilgisi almak için Refresh Token çözümle
	claims, err := utils.ParseAccessToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token geçersiz",
		})
	}

	email := claims["email"].(string)

	// 4. Redis'ten Refresh Token sil
	err = config.Rdb.Del(c.Context(), email).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "redis'ten token silinemedi",
		})
	}

	// 5. Cookie'yi sil
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout başarılı",
	})
}
