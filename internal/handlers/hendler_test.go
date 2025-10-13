package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	Key   string
	Value string
}

var theTests =[]struct {
	name   string
	url    string
	method string
	params []postData
	expectedStatusCode int

}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
		{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"sa", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-avail","/search-availability","POSt",[]postData{
	{Key:"start",Value: "2020-02-01"},
	{Key:"start",Value: "2020-02-02"},
},http.StatusOK},
	{"post-search-avail-json","/search-availability-json","POSt",[]postData{
	{Key:"start",Value: "2020-02-01"},
	{Key:"start",Value: "2020-02-02"},
},http.StatusOK},
	{"make reservation","/make-reservation","POSt",[]postData{
	{Key:"frist_name",Value: "John"},
	{Key:"Last-name",Value: "Smith"},
		{Key:"emil",Value: "m@ee.com"},
			{Key:"phone",Value: "555-55-5555"},
},http.StatusOK},
}
func TestHandlers(t *testing.T) {
	routes:= getRoutes()
	ts:=httptest.NewTLSServer(routes)
	defer ts.Close()

	for _,e:=range theTests{
		if e.method=="GET"{
		resp,err:=	ts.Client().Get(ts.URL+e.url)
		if err!=nil{
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode!=e.expectedStatusCode{
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}else{
		values:=url.Values{}
		for _,x :=  range e.params{
			values.Add(x.Key,x.Value)
		}
			resp,err:= ts.Client().PostForm(ts.URL + e.url , values)
			if err!=nil{
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode!=e.expectedStatusCode{
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

}
}
}