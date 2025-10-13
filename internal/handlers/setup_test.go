package handlers

import (
	"encoding/gob"
	"fmt"
	"os"

	"text/template"
	// "html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"

	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/models"
	"githup.com/Mo277210/booking/internal/render"
)


var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "../../templates"

var funcMap = template.FuncMap{}

func getRoutes() http.Handler {

	gob.Register(models.Reservation{})
app = config.AppConfig{}
	app.InProduction = false
	
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	
	}

	app.TemplateCache = tc
	app.UseCache = true
	//---------------------------------------------------------------------------------------------

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewTemplates(&app)
	
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//Developing our own middleware
	// mux.Use(WriteToConsole)
	//Creating handlers for our forms & adding CSRF Protection
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(Repo.Home))
	mux.Get("/about", http.HandlerFunc(Repo.About))
	mux.Get("/generals-quarters", http.HandlerFunc(Repo.Generals))
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", http.HandlerFunc(Repo.ReservationSummary))

	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", http.HandlerFunc(Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(Repo.AvailabilityJSON))


	mux.Get("/contact", http.HandlerFunc(Repo.Contact))
	//Enabling static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		Name:     "csrf_token",
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	dir, _ := os.Getwd()
log.Println("Current working directory:", dir)

	mycache := map[string]*template.Template{}

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