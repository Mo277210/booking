package main

//http://127.0.0.1:8085/about
//go mod tidy  عشان لو زوت حاجة في go.mod
//go mod init github.com/Mo2772/go_cors  عشان تعمل go mod init
//go get -u github.com/alexedwards/scs/v2  عشان تعمل go get
//استخدم الامر الذي في الاسفل لكي run الكود
//cd D:\visalstadio_code\heelo_ssswold\html-templates
//go run ./cmd/web/hello_world_web_app.go
//. Using pat for routing
// go run ./cmd/web/...
////////////////////////////////////////////////////////////////////////////
//go run hello_world_web_app.go hendlers.go render.go   <----------------- correct command to run the code
//////////////////////////////////////////////////////////////////////////
//frist attempt code (hello world web app)
// //second attempt code (Functions and handlers)
//3 attempt code (Reorganizing our code, and adding some basic styling to pages)
// 4th attempt code (Enabling Go Modules and refactoring our code to use packages)
//5th attempt code (Working with Layouts)

import (
	// "database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// "html/template"
	"github.com/alexedwards/scs/v2"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/handlers"
	"githup.com/Mo277210/booking/internal/helpers"
	"githup.com/Mo277210/booking/internal/models"
	"githup.com/Mo277210/booking/internal/render"
	"githup.com/Mo277210/booking/internal/driver"
)

const portNumber = ":8085"

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db,err := run()
	if err != nil {
		log.Fatal(err)
	}


	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("starting web server at port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB,error) {

	gob.Register(models.Reservation{})

	app.InProduction = false
	infoLog = log.New(os.Stdout, "INFP\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	
	//connection to the database
	log.Println("connecting to the database...")
	db, err :=driver.ConnectSQL("host=localhost port=5432 dbname=booking user=postgres password=1234")	
	if err !=nil{
		log.Fatal("can not connect to Dying")
	}
log.Println("connected to database")

	

tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache:", err)
		return nil,err
	}

	app.TemplateCache = tc
	app.UseCache = false
	//---------------------------------------------------------------------------------------------

	repo := handlers.NewRepo(&app,db)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
