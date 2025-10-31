package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"

	"testing"

	"githup.com/Mo277210/booking/internal/models"
)

type postData struct {
	Key   string
	Value string
}

var theTests =[]struct {
	name   string
	url    string
	method string
	expectedStatusCode int

}{
	{"home", "/", "GET",  http.StatusOK},
		{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	
// 	{"post-search-avail","/search-availability","POSt",[]postData{
// 	{Key:"start",Value: "2020-02-01"},
// 	{Key:"start",Value: "2020-02-02"},
// },http.StatusOK},
// 	{"post-search-avail-json","/search-availability-json","POSt",[]postData{
// 	{Key:"start",Value: "2020-02-01"},
// 	{Key:"start",Value: "2020-02-02"},
// },http.StatusOK},
// 	{"make reservation","/make-reservation","POSt",[]postData{
// 	{Key:"frist_name",Value: "John"},
// 	{Key:"Last-name",Value: "Smith"},
// 		{Key:"emil",Value: "m@ee.com"},
// 			{Key:"phone",Value: "555-55-5555"},
// },http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes:= getRoutes()
	ts:=httptest.NewTLSServer(routes)
	defer ts.Close()

	for _,e:=range theTests{
		resp,err:=	ts.Client().Get(ts.URL+e.url)
		if err!=nil{
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode!=e.expectedStatusCode{
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}


//TestRepostiory_Reservation test reservation 
func TestRepostiory_Reservation(t *testing.T) {

	reservation:= models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID: 1,
			RoomName: "General's Quarters",
		},
	}	
  req,_:=http.NewRequest("GET","/make-reservation",nil)

  ctx:=getCtx(req)
  req=req.WithContext(ctx)

  rr:=httptest.NewRecorder()
  reservation.RoomID=1

  session.Put(ctx,"reservation",reservation)

  handler:=http.HandlerFunc(Repo.Reservation)
  handler.ServeHTTP(rr,req)
	if rr.Code!=http.StatusOK{
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req,_=http.NewRequest("GET","/make-reservation",nil)
	ctx=getCtx(req)
	req=req.WithContext(ctx)

	rr=httptest.NewRecorder()
	handler.ServeHTTP(rr,req)
	if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusSeeOther)
	}
	// test with non-existent room
		req,_=http.NewRequest("GET","/make-reservation",nil)
	ctx=getCtx(req)
	req=req.WithContext(ctx)

	rr=httptest.NewRecorder()
}

func getCtx(req *http.Request) context.Context{
	ctx,err:=session.Load(req.Context(),req.Header.Get("X-Session"))
	if err!=nil{
		log.Println(err)
	}
	return ctx
}