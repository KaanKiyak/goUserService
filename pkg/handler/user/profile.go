package user

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"user-service/pkg/handler/user/model"
	"user-service/pkg/handler/utils"
)

// ProfileHandler kullanıcı profilini döner
// @Summary Kullanıcı Profil Bilgisi
// @Description Bearer token ile giriş yapan kullanıcının profil bilgilerini döner
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User "Kullanıcı bilgisi"
// @Failure 401 {object} ErrorResponse "Geçersiz veya eksik token"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /profile [get]
func ProfileHandler(c *fiber.Ctx) error {
	// 1. Authorization header kontrolü
	authHeader := c.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		model.NewEventLog(nil, "", "", "PROFILE_REQUEST", "FAILED", "authorization header missing or invalid", c.IP(), c.Get("User-Agent"), c.OriginalURL()).Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// 2. Token'ı ayır
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 3. Token parse et
	claims, err := model.ParseToken(tokenString)
	if err != nil {
		model.NewEventLog(nil, "", "", "PROFILE_REQUEST", "FAILED", "token parse failed", c.IP(), c.Get("User-Agent"), c.OriginalURL()).Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// 4. JWT doğrula
	if err := claims.ValidateJWT(string(utils.GetJWTSecret())); err != nil {
		model.NewEventLog(nil, "", "", "PROFILE_REQUEST", "FAILED", "token validation failed", c.IP(), c.Get("User-Agent"), c.OriginalURL()).Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// 5. user_id al
	userIDFloat, ok := claims.Payload["user_id"].(float64)
	if !ok {
		model.NewEventLog(nil, "", "", "PROFILE_REQUEST", "FAILED", "user_id missing in token", c.IP(), c.Get("User-Agent"), c.OriginalURL()).Save()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	userID := int(userIDFloat)

	// SessionID (uuid) güvenli al
	sessionID := "UNKNOWN"
	if val, ok := claims.Payload["uuid"].(string); ok {
		sessionID = val
	}

	// 6. Kullanıcı bilgisi DB'den çek
	user, err := utils.GetUserByID(userID)
	if err != nil {
		model.NewEventLog(&userID, "", sessionID, "PROFILE_REQUEST", "FAILED", "user not found in DB", c.IP(), c.Get("User-Agent"), c.OriginalURL()).Save()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid token"})
	}

	// Log SUCCESS
	model.NewEventLog(&userID, user.Email, sessionID, "PROFILE_REQUEST", "SUCCESS", "", c.IP(), c.Get("User-Agent"), c.OriginalURL()).Save()

	// 7. Kullanıcı bilgisi döndür
	return c.JSON(user)
}
