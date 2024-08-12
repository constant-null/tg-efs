package main

import (
	"net/url"

	"github.com/constant-null/tg-efs/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

var db *storage.Storage

//func checkHash(token) fiber.Handler {
//	return func(c fiber.Ctx) error {
//		c.Query()
//		return c.Next()
//	}
//}

func init() {
	db = storage.New()
}

func main() {
	app := fiber.New()
	//app.Use(checkHash)

	app.Use(cors.New(cors.Config{AllowOrigins: []string{"*"}}))
	app.Get("/menu", menu)
	app.Get("/sheet", menu)
	app.Get("/sheet_data", getSheetData)
	app.Post("/sheet_data", updateSheetData)

	app.Listen(":8080")
}

func menu(c fiber.Ctx) error {
	c.Set("Content-type", "text/html; charset=utf-8")

	return c.SendFile("sheet/sheet.min.html")
}

func updateSheetData(c fiber.Ctx) error {
	vals, _ := url.ParseQuery(string(c.Request().Header.Peek("X-User-Data")))
	userStr := vals.Get("user")
	var id = fastjson.GetInt([]byte(userStr), "id")
	var data map[string]interface{}
	if err := c.Bind().JSON(&data); err != nil {
		return errors.Wrap(err, "unable to decode json")
	}
	return db.Store(id, data)
}

func getSheetData(c fiber.Ctx) error {
	userStr := c.Query("user")
	var id = fastjson.GetInt([]byte(userStr), "id")
	data, err := db.Get(id)
	if err != nil {
		log.Error(err)
		return c.JSON(map[string]string{}, "application/json")
	}
	//data := map[string]string{
	//	"name":         "Васья",
	//	"background":   "Пилот",
	//	"class":        "Эльф",
	//	"organization": "Ассоциация Героев",
	//	"rank":         "1",
	//	"fame":         "0",
	//	"brutal":       "к4",
	//	"skillful":     "к6",
	//	"smart":        "к8",
	//	"charismatic":  "к10",
	//}
	return c.JSON(data, "application/json")
}
