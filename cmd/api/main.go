package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"log"
	_ "user-service/docs"
	"user-service/pkg/config"
	. "user-service/pkg/handler/user"
)

// @title  User Service API
// @version 1.0
// @description This is a sample website user service API for deneme.com.
// @termsOfService http://swagger.io/terms/

// @contact.name kaan Tech Team
// @contact.email tech@deneme.com
// @contact.url http://tech.deneme.com/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer {token}
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	config.MySqlConnect()
	config.RedisConnect()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/login", LoginHandler)
	app.Get("/logout", LogoutHandler)
	app.Get("/profile", ProfileHandler)
	app.Post("/register", RegisterHandler)
	app.Post("/refresh", RefreshHandle)
	err := app.Listen("localhost:8080")
	if err != nil {
		panic(err)
	}
}
