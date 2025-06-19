package handlers

import (
	"embed"
	"html/template"
)

type Handler struct {
	PublicFS embed.FS
	Version  string
	Base     *template.Template
}
