package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// تسجيل الـ struct المستخدم في الـ session
	gob.Register(models.Reservation{})

	// إعداد التطبيق
	testApp.InProduction = false

	infoLog:= log.New(os.Stdout,"INFP\t",log.Ldate|log.Ltime)
	testApp.InfoLog=infoLog

	errorLog:=log.New(os.Stdout,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog=errorLog
	// إعداد الـ session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
      
	testApp.Session = session
	app=&testApp

	os.Exit(m.Run())
}

type myWriter struct{


}

func (tw *myWriter) Header() http.Header {
	h := make(http.Header)
	return h
}


func (tw *myWriter)WriteHeader(i int){


}

func (tw *myWriter)Write(b[]byte)(int ,error){
	lenght:=len(b)
	return lenght,nil
}