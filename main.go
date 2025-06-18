package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Tool ToolSection `toml:"tool"`
}

type ToolSection struct {
	Commitizen CommitizenConfig `toml:"commitizen"`
}

type CommitizenConfig struct {
	Name                  string `toml:"name"`
	TagFormat             string `toml:"tag_format"`
	VersionScheme         string `toml:"version_scheme"`
	Version               string `toml:"version"`
	UpdateChangelogOnBump bool   `toml:"update_changelog_on_bump"`
	MajorVersionZero      bool   `toml:"major_version_zero"`
}

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

	e.FileFS("/favicon.ico", "public/favicon.ico", publicFS)

	e.GET("/up", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	f, err := os.ReadFile(".cz.toml")
	if err != nil {
		panic(err)
	}

	var config Config
	if _, err := toml.Decode(string(f), &config); err != nil {
		log.Fatalf("Failed to decode TOML: %v", err)
	}
	version := config.Tool.Commitizen.Version

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
			"title":        "UUID Generator",
			"uuid":         uuidValue,
			"selectedType": uuidType,
			"version":      version,
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
