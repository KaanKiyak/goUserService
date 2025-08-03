package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"user-service/pkg/config"
	"user-service/pkg/handler/user/model"
	"user-service/pkg/handler/utils"
)

// LogoutHandler kullanıcının oturumunu kapatır
// @Summary Kullanıcı Çıkışı (Logout)
// @Description Refresh token ve access token doğrulanır, Redis'ten refresh token ve session silinir, cookie temizlenir
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout başarılı"
// @Failure 401 {object} ErrorResponse "Token bulunamadı veya geçersiz"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /logout [post]
func LogoutHandler(c *fiber.Ctx) error {
	// 1. Cookie'den refresh token al
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		logEvent := model.NewEventLog(nil, "", "", "LOGOUT", "FAILED", "refresh token bulunamadı", c.IP(), c.Get("User-Agent"), c.OriginalURL())
		logEvent.Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh token bulunamadı"})
	}

	// 2. Authorization header'dan access token al
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		logEvent := model.NewEventLog(nil, "", "", "LOGOUT", "FAILED", "geçersiz token", c.IP(), c.Get("User-Agent"), c.OriginalURL())
		logEvent.Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "geçersiz token"})
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// 3. Access token parse et
	claimsAccess, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		logEvent := model.NewEventLog(nil, "", "", "LOGOUT", "FAILED", "access token geçersiz", c.IP(), c.Get("User-Agent"), c.OriginalURL())
		logEvent.Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "access token geçersiz"})
	}

	// 4. Kullanıcı bilgileri
	userID := int(claimsAccess["user_id"].(float64))

	//  sessionID güvenli al (fallback UNKNOWN)
	sessionID := "UNKNOWN"
	if val, ok := claimsAccess["uuid"].(string); ok {
		sessionID = val
	}

	// 5. Refresh token parse et
	claimsRefresh, _ := utils.ParseAccessToken(refreshToken)
	email := claimsRefresh["email"].(string)

	// 6. Redis'ten refresh token sil
	if err := config.Rdb.Del(c.Context(), email).Err(); err != nil {
		logEvent := model.NewEventLog(&userID, email, sessionID, "LOGOUT", "FAILED", "refresh token redis'ten silinemedi", c.IP(), c.Get("User-Agent"), c.OriginalURL())
		logEvent.Save()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "redis'ten token silinemedi"})
	}

	// 7. Redis'ten session sil
	redisKey := fmt.Sprintf("user:%d:session", userID)
	if err := config.Rdb.HDel(c.Context(), redisKey, sessionID).Err(); err != nil {
		logEvent := model.NewEventLog(&userID, email, sessionID, "LOGOUT", "FAILED", "session Redis'ten silinemedi", c.IP(), c.Get("User-Agent"), c.OriginalURL())
		logEvent.Save()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "session Redis'ten silinemedi"})
	}

	// 8. Cookie temizle
	utils.ClearSetRefreshToken(c)

	//  SUCCESS log
	logEvent := model.NewEventLog(&userID, email, sessionID, "LOGOUT", "SUCCESS", "", c.IP(), c.Get("User-Agent"), c.OriginalURL())
	logEvent.Save()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout başarılı"})
}
