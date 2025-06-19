package templater

import (
	"embed"
	"fmt"
	"html/template"

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

func (t *templater) FindTemplate(tmpl *template.Template, templateName string) error {
	_, err := tmpl.ParseFS(t.fs, fmt.Sprintf("public/%s.html", templateName))
	if err != nil {
		return echo.NewHTTPError(500, err)
	}
	return nil
}
