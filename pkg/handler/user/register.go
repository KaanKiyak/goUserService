package user

// RegisterHandler godoc
// @Summary Yeni kullanıcı kaydı
// @Description Kullanıcı bilgilerini veritabanına ekler.
// @Tags User
// @Accept json
// @Produce json
// @Param request body model.User true "Kayıt bilgileri"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"user-service/pkg/config"
	"user-service/pkg/handler/user/model"
)

// RegisterHandler godoc
func RegisterHandler(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Validasyonlar
	if err := user.ValidateUserName(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := user.ValidateEmail(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := user.ValidatePassword(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := user.ValidateRole(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := user.ValidateAge(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Kullanıcı UUID oluştur
	user.UUID = uuid.New().String()

	// DB'ye kaydet
	query := "INSERT INTO users (name, uuid, email, age, password, role) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := config.DB.Exec(query, user.Name, user.UUID, user.Email, user.Age, user.Password, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "MYSQL kayıt başarısız"})
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "ID alınamadı"})
	}
	user.ID = int(insertedID)

	// Redis'e kullanıcı bilgisi kaydet
	userJSON, _ := json.Marshal(user)
	redisKey := fmt.Sprintf("user:%d", user.ID)
	if err := config.Rdb.Set(c.Context(), redisKey, userJSON, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "kayıt başarılı",
		"user_id": user.ID,
		"uuid":    user.UUID,
	})
}
