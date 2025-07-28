package main

import (
	"embed"
	"html/template"
	"log"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/EwanGreer/uuid-generator/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Tool struct {
		Commitizen struct {
			Name                  string `toml:"name"`
			TagFormat             string `toml:"tag_format"`
			VersionScheme         string `toml:"version_scheme"`
			Version               string `toml:"version"`
			UpdateChangelogOnBump bool   `toml:"update_changelog_on_bump"`
			MajorVersionZero      bool   `toml:"major_version_zero"`
		} `toml:"commitizen"`
	} `toml:"tool"`
}

//go:embed public/*
var publicFS embed.FS

//go:embed .cz.toml
var czFile string

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Logger(), middleware.CORS())

	base := template.Must(template.New("base").ParseFS(publicFS, "public/base.html"))

	var config Config
	if _, err := toml.Decode(czFile, &config); err != nil {
		log.Fatalf("Failed to decode TOML: %v", err)
	}
	version := config.Tool.Commitizen.Version

	handler := handlers.Handler{
		Base:     base,
		Version:  version,
		PublicFS: publicFS,
	}

	e.StaticFS("/static", publicFS)
	e.FileFS("/favicon.ico", "public/favicon.ico", publicFS)

	e.GET("/up", handler.HandleUp)

	e.GET("/", handler.HandleIndexPage)
	e.GET("/password-generator", handler.HandlePasswordGeneratorPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
