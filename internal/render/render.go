package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/justinas/nosurf"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/models"
)

var funcMap = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"


//  NewRenderer returns a new template renderer

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	//
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	return td
}

//Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
    var tc map[string]*template.Template

    if app.UseCache {
        tc = app.TemplateCache
    } else {
        tc, _ = CreateTemplateCache()
    }

    t, ok := tc[fmt.Sprintf("%s.html", tmpl)]  // ✅ هنا التعديل
    if !ok {
        log.Println("❌ Template not found in cache:", tmpl)
        http.Error(w, "Template not found", 500)
        return errors.New("cannot get template from cache")
    }

    buf := new(bytes.Buffer)
    td = AddDefaultData(td, r)

    // ✅ إضافة فقط لتحديد القالب الرئيسي
    mainLayout := "base"
	
	// ✅ أي قالب يحتوي على كلمة "admin" يستخدم القالب الأساسي admin.html
	if strings.Contains(tmpl, "admin") {
		mainLayout = "admin"
	}

	err := t.ExecuteTemplate(buf, mainLayout, td)

    if err != nil {
        log.Println("❌ Error executing base template:", err)
        http.Error(w, "Internal Server Error", 500)
        return err
    }

    _, err = buf.WriteTo(w)
    if err != nil {
        log.Println("❌ Error writing template to browser:", err)
        http.Error(w, "Internal Server Error", 500)
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
		log.Println("📂 Looking for templates in:", pathToTemplates)
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

		// ✅ ربط القالب الأساسي base.html
		ts, err = ts.ParseFiles(fmt.Sprintf("%s/base.html", pathToTemplates))
		if err != nil {
			return mycache, err
		}

		// ✅ ربط قالب الأدمن admin.html إذا وُجد
		adminLayout := fmt.Sprintf("%s/admin.html", pathToTemplates)
		if _, err := os.Stat(adminLayout); err == nil {
			ts, _ = ts.ParseFiles(adminLayout)
		}

		mycache[name] = ts
	}

	return mycache, nil
}

