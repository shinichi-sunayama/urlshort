package app

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(dir string) *TemplateRenderer {
	pattern := filepath.Join(dir, "*.html")
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob(pattern)),
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
