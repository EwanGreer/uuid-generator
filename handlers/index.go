package handlers

import (
	"fmt"
	"html/template"
	"slices"

	"github.com/EwanGreer/uuid-generator/templater"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleIndexPage(c echo.Context) error {
	uuidType := c.QueryParam("type")

	if ok := slices.Contains([]string{"v4", "v7"}, uuidType); !ok {
		uuidType = "v4"
	}

	tmpl := template.Must(h.Base.Clone())
	t := templater.NewTemplater(h.PublicFS)

	err := t.FindTemplate(tmpl, "_nav", "index")
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
		"version":      h.Version,
	}

	return tmpl.ExecuteTemplate(c.Response(), "base", data)
}
