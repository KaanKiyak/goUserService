package user

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"user-service/pkg/handler/user/model"
	"user-service/pkg/handler/utils"
)

/*func ProfileHandler(c *fiber.Ctx) error {
	var user User
	authHeader := c.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseAccessToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	userID := int(claims["user_id"].(float64))

	err = config.DB.QueryRow("select id, email, role from users where id = ?", userID).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	return c.JSON(user)
}*/

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
	//authorization
	authHeader := c.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	// tokenı ayır
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// tokenı parse et
	claims, err := model.ParseToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	//validate
	if err := claims.ValidateJWT(string(utils.GetJWTSecret())); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	//payload expiration time çek

	//payloud user ıd çek
	userIDFloat, ok := claims.Payload["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	userID := int(userIDFloat)
	//kullanıcı bilgisini db den çek
	user, err := utils.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid token"})
	}

	return c.JSON(user)
}
