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
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"
// "html/template"
	"github.com/alexedwards/scs/v2"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/handlers"
	"githup.com/Mo277210/booking/internal/models"
	"githup.com/Mo277210/booking/internal/render"
)

const portNumber = ":8085"

var app config.AppConfig


func main() {

err:=run()
if err!=nil{
	log.Fatal(err)
}


	fmt.Println(fmt.Sprintf("starting web server at port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error{
	//what am i going to put in the session
   gob.Register(models.Reservation{})
	//7-----------(Setting application wide configuration)------------------------------------------------
	// change this to true when in production
	app.InProduction = false
	//---------------------------------------------------------------------------------------------
	//8-----------(Setting up a custom logger)------------------------------------------------

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	return err
	}

	app.TemplateCache = tc
	app.UseCache = false
	//---------------------------------------------------------------------------------------------

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return nil
}