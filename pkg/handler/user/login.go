package user

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"user-service/pkg/config"
	"user-service/pkg/handler/user/model"
	"user-service/pkg/handler/utils"
)

// ErrorResponse genel hata yanıtı için kullanılır
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LoginHandler kullanıcı girişi yapar
// @Summary Kullanıcı Girişi
// @Description Kullanıcı e-posta ve şifre ile giriş yapar, access ve refresh token döner
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.User true "Giriş bilgileri"
// @Success 200 {object} map[string]string "access_token döner"
// @Failure 400 {object} ErrorResponse "İstek hatalı"
// @Failure 401 {object} ErrorResponse "Geçersiz kullanıcı bilgisi"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /login [post]
func LoginHandler(c *fiber.Ctx) error {
	var req model.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Kullanıcı doğrulama
	var user model.User
	err := config.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil || user.Password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "invalid email or password",
		})
	}

	// Access & Refresh Token üret
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// Cookie Set + Redis Kaydet
	utils.SetRefreshToken(c, refreshToken)
	_ = config.Rdb.Set(c.Context(), user.Email, refreshToken, 24*time.Hour).Err()

	// Kullanıcıyı redisten takip et
	claims, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	userID := int(claims["user_id"].(float64))
	ip := c.IP()
	userAgent := c.Get("User-Agent")
	go model.TrackLogin(userID, ip, userAgent)

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}
