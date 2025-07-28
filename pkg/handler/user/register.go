package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"user-service/pkg/config"
	"user-service/pkg/handler/user/model"
)

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

func RegisterHandler(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	redisKey := fmt.Sprintf("user:%d", user.ID)
	if err := config.Rdb.Set(c.Context(), redisKey, string(c.Body()), 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := user.ValidateUserName()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = user.ValidateEmail()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"massage": err.Error(),
		})
	}
	err = user.ValidatePassword()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"massage": err.Error(),
		})
	}
	err = user.ValidateRole()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"massage": err.Error(),
		})
	}
	err = user.ValidateAge()
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"massage": err.Error(),
		})
	}

	/*
		if user.Age < 0 || user.Age >= 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "yaşama öl",
			})
		}
		if user.Role != "user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"massage": "bu sekme yanlızca user kaydı yapılır",
			})
		}
		for _, r := range user.Name {
			if !unicode.IsLetter(r) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "isim harflerden oluşmalıdır",
				})
			}
		}
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, user.Email)
		if !matched {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"massage": "düzgün formatta email gir",
			})
		}
	*/

	query := "INSERT INTO users (name, email, age, password, role) VALUES (?, ?, ?, ?, ?)"
	_, err = config.DB.Exec(query, user.Name, user.Email, user.Age, user.Password, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"massage": "MYSQL kayıt başarısız",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "kayıt başarılı",
	})
}
