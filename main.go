package main

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed public
var publicFS embed.FS

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Logger())

	e.StaticFS("/js", publicFS)
	base := template.Must(template.New("base").ParseFS(publicFS, "public/base.html"))

	e.GET("/", func(c echo.Context) error {
		uuidType := c.QueryParam("type")

		tmpl := template.Must(base.Clone())
		t := NewTemplater(publicFS)

		err := t.findTemplate(tmpl, "index")
		if err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Template parse error: %v", err))
		}

		var uid string
		switch uuidType {
		case "v7":
			id, err := uuid.NewV7()
			if err != nil {
				return err
			}
			uid = id.String()
		default:
			uuidType = "v4"
			uid = uuid.NewString()
		}

		data := map[string]any{
			"title":        "Page",
			"uuid":         uid,
			"selectedType": uuidType,
		}

		return tmpl.ExecuteTemplate(c.Response(), "base", data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}

type templater struct {
	fs embed.FS
}

func NewTemplater(fs embed.FS) *templater {
	return &templater{
		fs: fs,
	}
}

func (t *templater) findTemplate(tmpl *template.Template, templateName string) error {
	_, err := tmpl.ParseFS(t.fs, fmt.Sprintf("public/%s.html", templateName))
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Template parse error: %v", err))
	}
	return nil
}
