package handlers
// // الـ Handler = الدالة (function) أو الكائن (struct) اللي مسؤول عن معالجة الطلب (Request) والرد على العميل (Response).

// بمعنى آخر:

// الـ Route يحدد "المسار" (مثلاً /about).

// الـ Handler يحدد "إيه اللي هيحصل" لما المستخدم يدخل المسار ده
//Route = العنوان (URL).

// Handler = الكود اللي ينفذ لما نزور العنوان ده.

// تحب أديك مثال كامل صغير يجمع Route + Handler + Middleware عشان تبقى الصورة 100% 
//واضحة عندك؟
//كلمة Handlers (هاندلرز) معناها: المعالجين أو الوظائف المسؤولة عن التعامل مع الطلبات (Requests) في السيرفر.

// 🔹 في Go مع حزمة net/http:

// أي Handler هو دالة (أو Struct) بتنفذ شكل محدد:

// func(w http.ResponseWriter, r *http.Request)

// w http.ResponseWriter → اللي بيكتب فيه السيرفر الرد (Response).

// r *http.Request → بيمثل الطلب اللي جاي من المتصفح (URL, بيانات, Form...).

// كل Route (مسار) في الموقع بيرتبط بـ Handler.
// مثلًا:

// / → يروح لـ HomeHandler

// /about → يروح لـ AboutHandler
// Handlers = دوال بتتعامل مع طلبات HTTP.

// وظيفتها: تحدد إيه اللي يحصل لما المستخدم يزور URL معين (ترجع HTML, JSON, أو أي Response).

// غالبًا بتستخدم مع render.RenderTemplate عشان تعرض صفحات HTML.
import (
	"net/http"

	"githup.com/Mo277210/booking/pkg/config"
	"githup.com/Mo277210/booking/pkg/models"
	"githup.com/Mo277210/booking/pkg/render"
)

// في الـ Software Architecture (نمط Repository Pattern):

// Repository = طبقة وسيطة بين الكود (Business Logic) وقاعدة البيانات (Database).

// الهدف إنه يعزلك عن تفاصيل التعامل مع قاعدة البيانات.
// بدل ما تكتب SQL جوه الكود، تندّه على Repository.

////////////////////////////////////////////////////////////////////////////////////

//Repo the repository used by the handlers
var Repo *Respostory

//Respostory is the repository type
type Respostory struct{
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Respostory {
	return &Respostory{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Respostory) {
	Repo = r
}

//home handler
func (m *Respostory) Home(w http.ResponseWriter,r* http.Request){
remoteIP:=r.RemoteAddr
m.App.Session.Put(r.Context(),"remote_ip",remoteIP)

	render.RenderTemplate(w,r,"home",&models.TemplateData{})
}

//about handler

func (m *Respostory) About(w http.ResponseWriter,r* http.Request){
	//perform some logic
	stringMap:=make(map[string]string)
    stringMap["test"]="Hello, again. This is the about page"
  //Experimenting with sessions(video)
	remoteIP:=m.App.Session.GetString(r.Context(),"remote_ip")
   stringMap["remote_ip"]=remoteIP
	//send the data to the template
render.RenderTemplate(w,r,"about",&models.TemplateData{
	StringMap: stringMap,
})

}

// Reservation renders the make a reservation page and displays form
func (m *Respostory) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w,r, "make-reservation", &models.TemplateData{})
}

// Generals renders the room page
func (m *Respostory) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w,r, "generals", &models.TemplateData{})
}



// Majors renders the room page
func (m *Respostory) Majors(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w,r, "majors", &models.TemplateData{})
}


// Availability renders the search availability page
func (m *Respostory) Availability(w http.ResponseWriter, r *http.Request) {
 render.RenderTemplate(w,r, "search-availability", &models.TemplateData{})
}

// PostAvailability handles the post request of search availability form
func (m *Respostory) PostAvailability(w http.ResponseWriter, r *http.Request) {
  start:=r.Form.Get("start")
  end:=r.Form.Get("end")
  w.Write([]byte("start date is: "+start+" end date is: "+end))
	}


// Contact renders the contact page
func (m *Respostory) Contact(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r,"contacts", &models.TemplateData{})
}

