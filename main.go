package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role_id" default0:"1"`
}

type Article struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Filter  string `json:"filter"`
}

var DB *gorm.DB

func main() {
	app := fiber.New()
	app.Use(cors.New())
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost)/testgolang"), &gorm.Config{})

	if err != nil {
		panic("DB Cannot Connect")
	} else {
		fmt.Println("Connected")
	}

	app.Get("/", home)
	app.Post("/api/register", register)
	app.Post("/api/login", login)
	app.Post("/api/createArticle", createArticle)
	app.Get("/api/article", article)
	app.Delete("/api/article/:id", delArticle)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	DB = db

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Article{})
	app.Listen(":2000")

}

func home(c *fiber.Ctx) error {
	return c.SendString("Halo From Dio Blog")

}

func register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Context().Err()
	}
	user := User{
		Name:     data["name"],
		Email:    data["email"],
		Password: data["password"],
	}

	c.Accepts("application/json")
	DB.Create(&user)

	return c.JSON(&user)
}

func login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Fatal("Could Not Parsing This Data")
	}

	var user User

	DB.Where("email = ? AND password = ?", data["email"], data["password"]).First(&user)

	if user.Email == data["email"] && user.Password == data["password"] && user.Role == 0 || user.Role == 1 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "Login Success",
			"role":    user.Role,
		})
	}

	if user.Email != data["email"] && user.Password != data["password"] {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Login Failed",
		})
	}
	c.Accepts("application/json")
	return c.JSON(&user)
}

func createArticle(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {

		log.Fatal("Could Not Parsing This Data")
	}

	post := Article{
		Title:   data["title"],
		Content: data["content"],
		Filter:  data["filter"],
	}

	DB.Create(&post)

	return c.JSON(&post)

}

func article(c *fiber.Ctx) error {
	article := []Article{}

	DB.Find(&article)

	return c.JSON(&article)
}

func delArticle(c *fiber.Ctx) error {
	id := c.Params("id")

	var article Article

	DB.First(&article, id)
	DB.Delete(&article)

	return c.JSON(fiber.Map{
		"message": "Berhasil Dihapus Gan :v",
	})
}
