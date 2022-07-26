package render

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/tsawler/bookings/internal/config"
	"github.com/tsawler/bookings/internal/models"
)

// var app *config.AppConfig

// // NewTemplates sets the config for the template package
// func NewTemplates(a *config.AppConfig) {
// 	app = a
// }

// func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
// 	td.Flash = app.Session.PopString(r.Context(), "flash")
// 	td.Error = app.Session.PopString(r.Context(), "error")
// 	td.Warning = app.Session.PopString(r.Context(), "warning")
// 	td.CSRFToken = nosurf.Token(r)
// 	return td
// }

// //RenderTemplate renders HTML using template
// func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
// 	var tc map[string]*template.Template
// 	if app.UseCache {
// 		// get the template cache from the app config
// 		tc = app.TemplateCache
// 	} else {
// 		tc, _ = CreateTemplateCache()
// 	}

// 	// get the template cache from the app config

// 	// get requested template from cache
// 	t, ok := tc[tmpl]
// 	if !ok {
// 		log.Fatal("Could not get template from template cache")
// 	}

// 	buf := new(bytes.Buffer)

// 	td = AddDefaultData(td, r)
// 	_ = t.Execute(buf, td)

// 	// render the template
// 	_, err := buf.WriteTo(w)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	// parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
// 	// err := parsedTemplate.Execute(w, nil)
// 	// if err != nil {
// 	// 	fmt.Println("error parsing template:", err)
// 	// 	return
// 	// }
// }

// func CreateTemplateCache() (map[string]*template.Template, error) {
// 	myCache := map[string]*template.Template{}

// 	// get all of the files named *.page.tmpl from ./templates
// 	pages, err := filepath.Glob("./templates/*.page.tmpl")
// 	if err != nil {
// 		return myCache, err
// 	}

// 	// range through all files ending with *.page.tmpl
// 	for _, page := range pages {
// 		name := filepath.Base(page)
// 		ts, err := template.New(name).ParseFiles(page)
// 		if err != nil {
// 			return myCache, err
// 		}

// 		matches, err := filepath.Glob("./templates/*.layout.tmpl")
// 		if err != nil {
// 			return myCache, err
// 		}

// 		if len(matches) > 0 {
// 			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
// 			if err != nil {
// 				return myCache, err
// 			}
// 		}

// 		myCache[name] = ts
// 	}

// 	return myCache, nil

// }

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("Could not get template from template cache")
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
