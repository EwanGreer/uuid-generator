package main

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"slices"

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
	e.Use(middleware.Logger(), middleware.CORS())

	e.StaticFS("/static", publicFS)
	base := template.Must(template.New("base").ParseFS(publicFS, "public/base.html"))

	e.GET("/up", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/", func(c echo.Context) error {
		uuidType := c.QueryParam("type")

		if ok := slices.Contains([]string{"v4", "v7"}, uuidType); !ok {
			uuidType = "v4"
		}

		tmpl := template.Must(base.Clone())
		t := NewTemplater(publicFS)

		err := t.findTemplate(tmpl, "index")
		if err != nil {
			return echo.NewHTTPError(500, fmt.Sprintf("Template parse error: %v", err))
		}

		var uuidValue string
		switch uuidType {
		case "v7":
			id, err := uuid.NewV7()
			if err != nil {
				return err
			}
			uuidValue = id.String()
		default:
			uuidValue = uuid.NewString()
		}

		data := map[string]any{
			"title":        "Page",
			"uuid":         uuidValue,
			"selectedType": uuidType,
		}

		return tmpl.ExecuteTemplate(c.Response(), "base", data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(":" + port))
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
