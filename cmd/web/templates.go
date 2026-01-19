package main

import (
	"html/template"
	"io/fs"
	"modulo.porreiro/internal/models"
	"modulo.porreiro/ui"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	// value formatted according to the layout defined by the argument
	return t.Format("02 Feb 2020 at 16:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	// [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract the filename
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// parse the template into a template set
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)

		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	return cache, nil
}
