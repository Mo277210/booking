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
// cd "D:\visalstadio_code\heelo_ssswold\booking\MailHog"
// .\MailHog_windows_amd64.exe
//package main

//عشان تولد كلمة سر مشفرة
// import (
// 	"fmt"
// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	password := "password"

// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

// 	fmt.Println(string(hashedPassword))
// }

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
	"flag"
	"fmt"
	"log"
	"net/http"

	// "net/smtp"
	"os"
	"time"

	// "html/template"
	"github.com/alexedwards/scs/v2"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/driver"
	"githup.com/Mo277210/booking/internal/handlers"
	"githup.com/Mo277210/booking/internal/helpers"
	"githup.com/Mo277210/booking/internal/models"
	"githup.com/Mo277210/booking/internal/render"
)

const portNumber = ":8087"

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db,err := run()
	if err != nil {
		log.Fatal(err)
	}


	defer db.SQL.Close()
	defer close(app.MailChan)
	

// from := "me@here.com"
// auth := smtp.PlainAuth("", from, "1234", "localhost")
// err = smtp.SendMail("localhost:1025", auth, from, []string{"you@here.com"}, []byte("Subject: Test\r\n\r\nHello world"))
// if err != nil {
//     log.Println("Error sending email:", err)
// }



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
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})
	
	//read flages

inProduction:=flag.Bool("production",true,"Application is  in production" )
useCache:=flag.Bool("cache",true,"use template cache")
dbHost:=flag.String("dbhost","localhost","database host")
dbName:=flag.String("dbname","","database name")
dbUser:=flag.String("dbuser","","database user")
dbPassword:=flag.String("dbpass","","database password")
dbPort:=flag.String("dbport","5432","database port")
dbSSL:=flag.String("dbssl","disable","database ssl settings (disable, prefer, require)")

flag.Parse()

if *dbName=="" || *dbUser==""  {
	fmt.Println("missing required flags")
	os.Exit(1)
}


	mailchan := make(chan models.MailData)
	app.MailChan = mailchan
	
// شغّل مستمع البريد بعد ثانيتين من تشغيل السيرفر

	listenForMail()


	//---------------------------------------------------------------------------------------------

	app.InProduction = *inProduction
	app.UseCache = *useCache


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
	
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSSL)
	
	db, err :=driver.ConnectSQL(connectionString)	
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

	//---------------------------------------------------------------------------------------------

	repo := handlers.NewRepo(&app,db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
