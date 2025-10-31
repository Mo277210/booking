package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

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
//TestRepository_PostReservation tests PostReservation handler
func TestRepository_PostReservation(t *testing.T){

	reqBody:="start_date=2050-01-01"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-10")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=John")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=1")


  req,_:=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))

  ctx:=getCtx(req)
  req=req.WithContext(ctx)

  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr:=httptest.NewRecorder()

  handler:=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusSeeOther)
	}
	// test case where first name is too short
  req,_=http.NewRequest("POST","/make-reservation",nil)
 ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code for missing body: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}
	//test for invalid start date

	reqBody="start_date=invalid"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-10")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=John")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=1")

  req,_=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
  ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

		//test for invalid end date

	reqBody="start_date=2050-01-01"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=invalid")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=John")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=1")

  req,_=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
  ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

		//test for invalid room id

	reqBody="start_date=2050-01-01"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-10")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=John")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=invalid")

  req,_=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
  ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}
//test for invalid data in form fields (first name too short)
	reqBody="start_date=2050-01-01"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-10")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=j")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=1")

  req,_=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
  ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

	//test for failure to insert reservation into database
	reqBody="start_date=2050-01-01"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-10")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=jhon")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=2")

  req,_=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
  ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler failed when trying to insert reservation: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

		//test for failure to insert restriction into database
	reqBody="start_date=2050-01-01"
	reqBody=fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-10")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"first_name=jhon")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"email=m@ee.com")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"phone=555-555-5555")
	reqBody=fmt.Sprintf("%s&%s",reqBody,"room_id=1000")

  req,_=http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
  ctx=getCtx(req)
  req=req.WithContext(ctx)
  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
  rr=httptest.NewRecorder()

  handler=http.HandlerFunc(Repo.PostReservation)
  handler.ServeHTTP(rr,req)

if rr.Code!=http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler failed when trying to insert reservation: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}
		}

func TestRepository_AvailabilityJSON(t *testing.T){

	//first  case -rooms are not available
	resBody:="start=2050-01-01"
	resBody=fmt.Sprintf("%s&%s",resBody,"end=2050-01-10")
	resBody=fmt.Sprintf("%s&%s",resBody,"room_id=1")

	
	//create request
	  req,_:=http.NewRequest("POST","/search-availability-json",strings.NewReader(resBody))
	  
	  //get context from request
	  ctx:=getCtx(req)
	  req=req.WithContext(ctx)
	  
	  //set the request header
	  req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	 
	  //make handler handlerfunc
	  handler:=http.HandlerFunc(Repo.AvailabilityJSON)
	
	  //got response recorder
	  rr:=httptest.NewRecorder()
	
	  //serve the http
	  handler.ServeHTTP(rr,req)

	  var j jsonResponse

	  err:=json.Unmarshal([]byte(rr.Body.String()),&j)
	  if err!=nil{
		  t.Error("failed to parse json")
	  }



}

func getCtx(req *http.Request) context.Context{
	ctx,err:=session.Load(req.Context(),req.Header.Get("X-Session"))
	if err!=nil{
		log.Println(err)
	}
	return ctx
}

// | العملية                | الدالة        | من             | إلى            | مثال                                  |
// | ---------------------- | ------------- | -------------- | -------------- | ------------------------------------- |
// | **تحليل النص إلى وقت** | `time.Parse`  | `"2025-10-26"` | `time.Time`    | إدخال من المستخدم                     |
// | **تحويل الوقت إلى نص** | `time.Format` | `time.Time`    | `"2025-10-26"` | عرض للمستخدم أو حفظ في قاعدة البيانات |
