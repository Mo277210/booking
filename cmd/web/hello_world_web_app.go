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
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/handlers"
	"githup.com/Mo277210/booking/internal/render"
)

const portNumber = ":8085"

var app config.AppConfig

// var session *scs.SessionManager  // ❌ احذف ده، مش محتاجه
// var infoLog *log.Logger
// var errorLog *log.Logger
// main function
func main() {
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
	}

	app.TemplateCache = tc
	app.UseCache = false
	//---------------------------------------------------------------------------------------------

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//Using pat for routing

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("starting web server at port %s", portNumber))
	//Using pat for routing
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

/////3 attempt code (Reorganizing our code, and adding some basic styling to pages)////////////////////////////////////////////////////////////////
// package main

// //second attempt code (Serving HTML Templates)
// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"text/template"
// )
// const portNumber=":8080"

// //home handler
// func Home(w http.ResponseWriter,r* http.Request){
// renderTemplate(w,"home")
// }

// //about handler

// func About(w http.ResponseWriter,r* http.Request){
// renderTemplate(w,"about")

// }

// func renderTemplate(w http.ResponseWriter, tmpl string) {
//     path := "./html-templates/" + tmpl + ".html"
//     if _, err := os.Stat(path); os.IsNotExist(err) {
//         http.Error(w, "Template file not found: "+path, http.StatusInternalServerError)
//         return
//     }
//     parsedTemplate, err := template.ParseFiles(path)
//     if err != nil {
//         http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     err = parsedTemplate.Execute(w, nil)
//     if err != nil {
//         http.Error(w, "Error rendering template", http.StatusInternalServerError)
//         return
//     }
// }

// //main function
// func main() {

// http.HandleFunc("/",Home)
// http.HandleFunc("/about",About)
// fmt.Println(fmt.Sprintf("starting web server at port %s",portNumber))
// 	_=http.ListenAndServe(portNumber,nil)
// }

/////////////11111111111111111111111111111111111111111111111111111111//////////////////////////////////////////////////////////////////////
//frist attempt code (hello world web app)

////package main

// import (
// 	"fmt"
// 	"net/http"
// )

// func main() {

// 	http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){

// 		n, err:=fmt.Fprintf(w,"Hello, World!")

// 		if err!=nil {
// 			fmt.Println("Error writing response:",err)
// 			return
// 		}

// 		fmt.Println(fmt.Sprintf("n=%d, err=%v",n,err))

// 	})

// 	_=http.ListenAndServe(":8080",nil)
// }

////////////////////////22222222222222222222222222222222////////////////////////////////////////////////////////////////////////////////////////

// package main

// //second attempt code (Functions and handlers)
// import (
// 	"fmt"
// 	"net/http"
// )
// const portNumber=":8080"

// //home handler
// func Home(w http.ResponseWriter,r* http.Request){
// fmt.Fprintf(w,"this is home page")

// }

// //about handler

// func About(w http.ResponseWriter,r* http.Request){
// 	sum:=AddValue(5,10)
// _,_= fmt.Fprintf(w,"the sum is %d\n",sum)

// }
// //function to add two values
// func AddValue(x int,y int)int{
// return x+y
// }

// func Divide(w http.ResponseWriter,r* http.Request){

// 	f,err:=divideValues(100.0,0.0)
// 	if err!=nil{
// 		fmt.Fprintf(w,"Cannot divide by zero")
// 		return
// 	}
// fmt.Fprintf(w,fmt.Sprintf("%f divided by %f is %f",100.0,0.0,f))
// }
// func divideValues(x,y float64)(float64,error){
// 	if y<=0{
// 		err:=fmt.Errorf("cannot divide by zero  ")
// 		return 0,err
// 	}
// 	result:=x/y
// 	return result,nil
// }

// //main function
// func main() {

// http.HandleFunc("/",Home)
// http.HandleFunc("/about",About)
// http.HandleFunc("/divide",Divide)
// fmt.Println(fmt.Sprintf("starting web server at port %s",portNumber))
// 	_=http.ListenAndServe(portNumber,nil)
// }
