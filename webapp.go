package main

import (
	_ "embed"
	"net/url"
	"os"

	"github.com/constant-null/tg-efs/storage"
	"github.com/fxamacker/cbor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

var db storage.Storage

//go:embed sheet/sheet.min.html
var sheetHtml string

type SheetData struct {
	Name         string `json:"name" cbor:"1,keyasint,omitempty"`
	Background   string `json:"background" cbor:"2,keyasint,omitempty"`
	Class        string `json:"class" cbor:"3,keyasint,omitempty"`
	Organization string `json:"organization" cbor:"4,keyasint,omitempty"`
	Rank         string `json:"rank" cbor:"5,keyasint,omitempty"`
	Fame         string `json:"fame" cbor:"6,keyasint,omitempty"`
	Brutal       string `json:"brutal" cbor:"7,keyasint,omitempty"`
	Skillful     string `json:"skillful" cbor:"8,keyasint,omitempty"`
	Smart        string `json:"smart" cbor:"9,keyasint,omitempty"`
	Charismatic  string `json:"charismatic" cbor:"10,keyasint,omitempty"`
}

var DefaultSheetData = SheetData{
	Name:         "",
	Background:   "",
	Class:        "",
	Organization: "",
	Rank:         "1",
	Fame:         "0",
	Brutal:       "к4",
	Skillful:     "к4",
	Smart:        "к4",
	Charismatic:  "к4",
}

//func checkHash(token) fiber.Handler {
//	return func(c fiber.Ctx) error {
//		c.Query()
//		return c.Next()
//	}
//}

func init() {
	var err error
	if os.Getenv("DEBUG") != "" {
		db = storage.NewLocal()
	} else {
		db, err = storage.New()
		if err != nil {
			log.Fatalf("Initialization error: %+v", err)
		}
	}
}

func runWebApp() {
	app := fiber.New()
	//app.Use(checkHash)

	app.Use(cors.New(cors.Config{AllowOrigins: []string{"*"}}))
	app.Get("/menu", getSheet)
	app.Get("/sheet", getSheet)
	app.Get("/sheet_data", getSheetData)
	app.Post("/sheet_data", updateSheetData)

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func getSheet(c fiber.Ctx) error {
	c.Set("Content-type", "text/html; charset=utf-8")

	return c.SendString(sheetHtml)
}

func updateSheetData(c fiber.Ctx) error {
	vals, _ := url.ParseQuery(string(c.Request().Header.Peek("X-User-Data")))
	userStr := vals.Get("user")
	var id = int64(fastjson.GetInt([]byte(userStr), "id"))
	if id == 0 {
		c.Status(400)
		c.WriteString("missing id")
		return nil
	}

	var data SheetData
	if err := c.Bind().JSON(&data); err != nil {
		return errors.Wrap(err, "unable to decode json")
	}
	bdata, _ := cbor.Marshal(data, cbor.EncOptions{})
	return db.Store(c.Context(), id, bdata)
}

func getSheetData(c fiber.Ctx) error {
	userStr := c.Query("user")
	var id = int64(fastjson.GetInt([]byte(userStr), "id"))
	if id == 0 {
		c.Status(400)
		c.WriteString("missing id")
		return nil
	}

	var sheetData SheetData
	data, err := db.Get(c.Context(), id)
	if err != nil {
		log.Error(err)
		return c.JSON(DefaultSheetData, "application/json")
	}

	if err := cbor.Unmarshal(data, &sheetData); err != nil {
		log.Error(err)
		return c.JSON(DefaultSheetData, "application/json")
	}

	return c.JSON(sheetData, "application/json")
}
