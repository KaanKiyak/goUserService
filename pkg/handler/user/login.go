package user

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
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
		// Log başarısız parse (optional)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// DB'den kullanıcıyı bul
	var user model.User
	err := config.DB.QueryRow("SELECT id, uuid, email, password FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.UUID, &user.Email, &user.Password)

	// Kullanıcı yok veya şifre hatalı
	if err != nil || user.Password != req.Password {
		//  Log FAILED LOGIN
		eventLog := model.NewEventLog(nil, req.Email, "", "LOGIN", "FAILED", "invalid email or password", c.IP(), c.Get("User-Agent"), c.OriginalURL())
		eventLog.Save()

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid email or password"})
	}

	//  Yeni session ID üret
	sessionID := uuid.New().String()

	//  JWT oluştur (uuid = session ID)
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	//  Refresh token cookie set
	utils.SetRefreshToken(c, refreshToken)

	//  Refresh token Redis'e kaydet (email -> token)
	_ = config.Rdb.Set(c.Context(), user.Email, refreshToken, 24*time.Hour).Err()

	//  Session bilgisi Redis HASH'e kaydet (user:<id>:session)
	sessionData := map[string]string{
		"ip":        c.IP(),
		"userAgent": c.Get("User-Agent"),
		"loginTime": time.Now().Format(time.RFC3339),
	}
	sessionJSON, _ := json.Marshal(sessionData)
	redisKey := fmt.Sprintf("user:%d:session", user.ID)
	if err := config.Rdb.HSet(c.Context(), redisKey, sessionID, sessionJSON).Err(); err != nil {
		log.Println("Redis session kaydı hatası:", err)
	}

	//  Log SUCCESS LOGIN
	userID := user.ID
	eventLog := model.NewEventLog(&userID, user.Email, sessionID, "LOGIN", "SUCCESS", "", c.IP(), c.Get("User-Agent"), c.OriginalURL())
	eventLog.Save()

	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"session_id":   sessionID,
	})
}
