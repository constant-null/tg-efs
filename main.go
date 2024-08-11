package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fastjson"
)

//func checkHash(token) fiber.Handler {
//	return func(c fiber.Ctx) error {
//		c.Query()
//		return c.Next()
//	}
//}

func main() {
	app := fiber.New()
	//app.Use(checkHash)
	// GET: http://localhost:8080/menu
	app.Get("/menu", menu)
	app.Get("/sheet", menu)
	app.Get("/test", test)

	app.Listen(":8080")
}

func menu(c fiber.Ctx) error {
	c.Set("Content-type", "text/html; charset=utf-8")
	c.WriteString(menuTpl)
	return nil
}

func test(c fiber.Ctx) error {
	userStr := c.Query("user")
	var id = fastjson.GetInt([]byte(userStr), "id")
	fmt.Println(id)
	c.Set("Content-type", "text/html; charset=utf-8")
	return nil
}
