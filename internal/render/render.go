package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/models"
)

var funcMap = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"



func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	//
	td.CSRFToken = nosurf.Token(r)

	return td
}


func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error{
	
	var tc map[string]*template.Template

	if app.UseCache {
	
		tc = app.TemplateCache
	} else {

		tc, _ = CreateTemplateCache()

	}


	t, ok := tc[tmpl]
	if !ok {
		return errors.New(	"canot get template from cache!")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	err := t.ExecuteTemplate(buf, "base", td)

	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("errr writing template to borwser",err)
		return err
	}
	return nil
	
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	mycache := map[string]*template.Template{}
	cwd, _ := filepath.Abs(".")
	log.Println("🔍 Current working directory:", cwd)
	log.Println("🔍 Template path being searched:", pathToTemplates)


	pages, err := filepath.Glob(fmt.Sprintf("%s/*.html", pathToTemplates))
	if err != nil {
		return mycache, err
	}

	log.Println("Looking for templates in ./templates/*.html")
	log.Println("Templates found:", pages)

	for _, page := range pages {
		name := filepath.Base(page)
		log.Println("Parsing template:", name)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return mycache, err
		}

		ts, err = ts.ParseFiles(fmt.Sprintf("%s/base.html", pathToTemplates))
		if err != nil {
			return mycache, err
		}

		mycache[name] = ts
	}

	return mycache, nil
}

