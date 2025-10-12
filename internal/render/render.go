package render

// كلمة Render (رندر) معناها الحرفي "عرض" أو "تجسيد".
// في مجال الويب والتطبيقات، معناها: تحويل القالب (Template) + البيانات → لصفحة HTML تُرسل للمستخدم.
import (
	"bytes"
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

// NewTemplates sets the config for the template package

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefault(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)

	return td
}

//Building a more complex template cache
// نسخة للتجارب السريعة (من غير كاش)
//Optimizing our template cache by using an application config

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// if we are using the template cache
	var tc map[string]*template.Template

	if app.UseCache {
		// use the template cache

		//get the template cache from app config
		tc = app.TemplateCache
	} else {

		tc, _ = CreateTemplateCache()

	}

	//create template cache

	// tc, err := CreateTemplateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//get requested template from cache
	t, ok := tc[tmpl+".html"]
	if !ok {
		http.Error(w, "Template not found in cache: "+tmpl, http.StatusInternalServerError)
		return
	}
	//render the template
	buf := new(bytes.Buffer)
	td = AddDefault(td, r)
	_ = t.Execute(buf, td)

	// ننفذ باستخدام "base" عشان القالب الأساسي يشتغل
	err := t.ExecuteTemplate(buf, "base", td)

	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	////befor (قبل) Building a simple template cache
	// path := "templates/" + tmpl + ".html"
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	http.Error(w, "Template file not found: "+path, http.StatusInternalServerError)
	// 	return
	// }

	// parsedTemplate, err := template.ParseFiles(
	// 	path,                  // القالب المطلوب (home.html أو about.html)
	// 	"templates/base.html", // القالب الأساسي (مهم: لازم ييجي بعد الصفحة)
	// )
	// if err != nil {
	// 	http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// err = parsedTemplate.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	http.Error(w, "Error rendering template", http.StatusInternalServerError)
	// 	return
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	mycache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.html")
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

		ts, err = ts.ParseFiles("./templates/base.html")
		if err != nil {
			return mycache, err
		}

		mycache[name] = ts
	}

	return mycache, nil
}

//func CreateTemplateCache() (map[string]*template.Template, error) {
// 	mycache := map[string]*template.Template{}

// 	// نقرأ كل ملفات الصفحات (home.html, about.html ...)
// 	pages, err := filepath.Glob("./html-templates/templates/*.html")

// 	if err != nil {
// 		return mycache, err
// 	}

// 	// نعمل لوب على كل صفحة
// 	for _, page := range pages {
// 		name := filepath.Base(page)

// 		// أول حاجة: نبدأ بالصفحة المطلوبة
// 		ts, err := template.ParseFiles(page)
// 		log.Println("Loading template into cache:", name)
// 		if err != nil {
// 			return mycache, err
// 		}

// 		// بعدين نضيف الـ base.html (مهم: ييجي بعد الصفحة)
// 		ts, err = ts.ParseFiles("./html-templates/templates/base.html")
// 		if err != nil {
// 			return mycache, err
// 		}

// 		// نضيف التيمبلت للكاش
// 		mycache[name] = ts
// 	}

// 	return mycache, nil
// }

//Building a simple template cache
// كاش للتيمبلت
// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	// check cache
// 	_, inMap := tc[t]
// 	if !inMap {
// 		log.Println("Creating template and adding to cache...")
// 		err = CreateTemplateCache(t)
// 		if err != nil {
// 			log.Println("Error creating template cache:", err)
// 			http.Error(w, "Error creating template cache", http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		log.Println("Using cached template:", t)
// 	}

// 	// get template from cache
// 	tmpl = tc[t]

// 	// لازم ننفذ باسم القالب الأساسي
// 	err = tmpl.ExecuteTemplate(w, "base", nil)
// 	if err != nil {
// 		log.Println("Error executing template:", err)
// 		http.Error(w, "Error rendering template", http.StatusInternalServerError)
// 		return
// 	}
// }

// func CreateTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("templates/%s.html", t),
// 		"templates/base.html",
// 	}

// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	// add template to cache (map)
// 	tc[t] = tmpl
// 	return nil
// }
