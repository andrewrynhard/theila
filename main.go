package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/talos-systems/theila/pkg/render"
	"github.com/talos-systems/theila/pkg/routes"
)

var (
	//go:embed frontend/dist
	dist embed.FS

	//go:embed frontend/src/templates
	templates embed.FS

	t *template.Template
)

func init() {
	var layoutFound bool

	err := fs.WalkDir(templates, "frontend/src/templates", func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsRegular() {
			return nil
		}

		name := filepath.Base(path)

		b, err := templates.ReadFile(path)
		if err != nil {
			return nil
		}

		s := string(b)

		var tmpl *template.Template

		if t == nil {
			t = template.New(name)
		}

		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		log.Printf("Parsing template %q", path)

		_, err = tmpl.Parse(s)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, tmpl := range t.Templates() {
		log.Printf("Defined template %q", tmpl.Name())
		if tmpl.Name() == render.LayoutTemplateName {
			layoutFound = true
		}
	}

	if !layoutFound {
		panic("a layout template is required")
	}
}

func main() {
	log.Println("Starting app")

	router := registerRoutes()

	port := 8080

	log.Printf("Serving on port %d", port)

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes() http.Handler {
	log.Println("Registering routes")

	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	renderer := &render.Renderer{Template: t}

	routes.Register(r, renderer)

	r.Handle("/frontend/dist/*", http.FileServer(http.FS(dist)))

	logRoutes(r)

	return r
}

func logRoutes(r chi.Router) {
	log.Println("Serving with routes:")
	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("\t%s %s", method, route)
		return nil
	})
}
