package templater

import (
	"embed"
	"html/template"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type templater struct {
	fs embed.FS
}

func NewTemplater(fs embed.FS) *templater {
	return &templater{
		fs: fs,
	}
}

func (t *templater) FindTemplate(tmpl *template.Template, templateNames ...string) error {
	joinedPaths := []string{}
	for _, name := range templateNames {
		name = name + ".html"
		joinedPaths = append(joinedPaths, filepath.Join("public", name))
	}

	_, err := tmpl.ParseFS(t.fs, joinedPaths...)
	if err != nil {
		return echo.NewHTTPError(500, err)
	}
	return nil
}
