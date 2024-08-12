package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/url"
	"os"

	"github.com/constant-null/tg-efs/storage"
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
	Name         string `json:"name"`
	Background   string `json:"background"`
	Class        string `json:"class"`
	Organization string `json:"organization"`
	Rank         string `json:"rank"`
	Fame         string `json:"fame"`
	Brutal       string `json:"brutal"`
	Skillful     string `json:"skillful"`
	Smart        string `json:"smart"`
	Charismatic  string `json:"charismatic"`
}

func (s *SheetData) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(
		b,
		&s.Name,
		&s.Background,
		&s.Class,
		&s.Organization,
		&s.Rank,
		&s.Fame,
		&s.Brutal,
		&s.Skillful,
		&s.Smart,
		&s.Charismatic,
	)

	return err
}

func (s *SheetData) MarshalBinary() (data []byte, err error) {
	var b bytes.Buffer
	fmt.Fprintln(
		&b,
		s.Name,
		s.Background,
		s.Class,
		s.Organization,
		s.Rank,
		s.Fame,
		s.Brutal,
		s.Skillful,
		s.Smart,
		s.Charismatic,
	)
	return b.Bytes(), nil
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

func main() {
	app := fiber.New()
	//app.Use(checkHash)

	app.Use(cors.New(cors.Config{AllowOrigins: []string{"*"}}))
	app.Get("/,enu", getSheet)
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
	bdata, _ := data.MarshalBinary()
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
		return c.JSON(sheetData, "application/json")
	}
	if err := sheetData.UnmarshalBinary(data); err != nil {
		log.Error(err)
		return c.JSON(sheetData, "application/json")
	}

	return c.JSON(sheetData, "application/json")
}
