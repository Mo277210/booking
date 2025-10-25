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

// غالبًا بتستخدم مع render.Template عشان تعرض صفحات HTML.
import (
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/driver"
	"githup.com/Mo277210/booking/internal/forms"
	"githup.com/Mo277210/booking/internal/helpers"
	"githup.com/Mo277210/booking/internal/models"
	"githup.com/Mo277210/booking/internal/render"
	"githup.com/Mo277210/booking/internal/repostiory"
	"githup.com/Mo277210/booking/internal/repostiory/dbrepo"
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
    DB repostiory.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig,db *driver.DB) *Respostory {
	return &Respostory{
		App: a,
        DB: dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Respostory) {
	Repo = r
}

//home handler
func (m *Respostory) Home(w http.ResponseWriter,r* http.Request){

	render.Template(w,r,"home",&models.TemplateData{})
}

//about handler

func (m *Respostory) About(w http.ResponseWriter,r* http.Request){

	
	//send the data to the template
render.Template(w,r,"about",&models.TemplateData{})

}

// Reservation renders the make a reservation page and displays form
func (m *Respostory) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
    data := make(map[string]interface{})
    data["reservation"] = emptyReservation

    render.Template(w,r, "make-reservation", &models.TemplateData{
        Form: forms.New(nil),
        Data: data,
    })
}

// PostReservation handles the posting of a reservation form
func (m *Respostory) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}


    sd:=r.Form.Get("start_date")
    ed:=r.Form.Get("end_date")

    //2020-01-01 --01/02 03:04:05PM '06-0700
    layout:="2006-01-02"
    StartDate,err:=time.Parse(layout,sd)
    if err!=nil {
        helpers.ServerError(w,err)
        return
    }
    endDate,err:=time.Parse(layout,ed)
    if err!=nil {
        helpers.ServerError(w,err)
        return
    }

    roomID,err:=strconv.Atoi(r.Form.Get("room_id"))
    if err!=nil {
        helpers.ServerError(w,err)
        return
    }

  
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
        StartDate:StartDate ,
        EndDate:endDate,
        RoomID:    roomID,
        
	}

	form := forms.New(r.PostForm)
	// form.Has("first_name", r)
    form.Required("first_name", "last_name", "email")
    form.MinLength("first_name", 3)
    form.IsEmail("email")


	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

  newReservationID,err:=  m.DB.InsertReservation(reservation)
  if err!=nil {
      helpers.ServerError(w,err)
      return
  }

  restriction:=models.RoomRestrictions{

	StartDate    : StartDate,
	EndDate      : endDate,
	RoomID       : roomID,
	ReservationID: newReservationID,
	RestrictionID: 1,
  }

err=m.DB.InsertRoomRestriction(restriction)
if err!=nil {
    helpers.ServerError(w,err)
    return
}
 

    m.App.Session.Put(r.Context(), "reservation", reservation)
    http.Redirect(w,r,"/reservation-summary",http.StatusSeeOther)
//if we reach here means the form is valid so we can put the reservation in the session ويحمل الصحفة من تاني
	// http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}






// Generals renders the room page
func (m *Respostory) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w,r, "generals", &models.TemplateData{})
}



// Majors renders the room page
func (m *Respostory) Majors(w http.ResponseWriter, r *http.Request) {
    render.Template(w,r, "majors", &models.TemplateData{})
}


// Availability renders the search availability page
func (m *Respostory) Availability(w http.ResponseWriter, r *http.Request) {
 render.Template(w,r, "search-availability", &models.TemplateData{})
}

// PostAvailability handles the post request of search availability form
func (m *Respostory) PostAvailability(w http.ResponseWriter, r *http.Request) {
  start:=r.Form.Get("start")
  end:=r.Form.Get("end")

        layout:="2006-01-02"
    StartDate,err:=time.Parse(layout,start)
    if err!=nil {
        helpers.ServerError(w,err)
        return
    }
    endDate,err:=time.Parse(layout,end)
    if err!=nil {
        helpers.ServerError(w,err)
        return
    }

    rooms,err:=m.DB.SearchAvailabilityForAllRooms(StartDate,endDate)

    if err!=nil {
        helpers.ServerError(w,err)
        return
    }


    if len(rooms)==0{
    m.App.Session.Put(r.Context(), "error", "No availability")
    http.Redirect(w,r,"/search-availability",http.StatusSeeOther)
    return
    }

    data:=make(map[string]interface{})
    data["rooms"]=rooms
// res
res := models.Reservation{
    StartDate: StartDate, 
    EndDate:   endDate,
  
}

m.App.Session.Put(r.Context(), "reservation", res)

  render.Template(w, r, "choose-room", &models.TemplateData{
        Data: data,
    })

	}

	type jsonResponse struct{
		OK bool `json:"ok"`
		Message string `json:"message"`
	}
// AvailabilityJSON handles request for availability and sends JSON response
func (m *Respostory) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    resp := struct {
        OK      bool   `json:"ok"`
        Message string `json:"message"`
    }{
        OK:      true,
        Message: "Available!",
    }

    out, err := json.Marshal(resp)
    if err != nil {
        helpers.ServerError(w,err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(out)
}


// Contact renders the contact page
func (m *Respostory) Contact(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r,"contacts", &models.TemplateData{})
}


func (m *Respostory) ReservationSummary(w http.ResponseWriter, r *http.Request) {
   
    reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
    if !ok {
       m.App.ErrorLog.Println("canot get error from session")
        m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
       
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }
    //remove reservation from session
    // because we don't need it anymore
    // and to avoid showing the same data again if user refresh the page
    // or go back to the summary page again
    m.App.Session.Remove(r.Context(), "reservation")

    data := make(map[string]interface{})
    data["reservation"] = reservation

    m.App.Session.Remove(r.Context(), "reservation")

    render.Template(w, r, "reservation-summary", &models.TemplateData{
        Data: data,
    })
}



// //ممتاز جدًا 🙌
// أنت الآن داخل مشروع Go (Golang) يستخدم **net/http** و **chi router**، وسؤالك عن `header` في هذا السياق ذكي جدًا 👏
// خلينا نشرحها بدقة وبأسلوب عملي مع أمثلة من كودك.

// ---

// ## 🧩 أولاً: ما هو الـ **Header** في HTTP؟

// الـ **HTTP Header** هو جزء من الطلب (Request) أو الرد (Response) يُستخدم لإرسال **معلومات إضافية (Metadata)** بين العميل (المتصفح أو التطبيق) والسيرفر.

// * الـ **Request Header** → يرسله المتصفح إلى السيرفر.
// * الـ **Response Header** → يرسله السيرفر إلى المتصفح.

// ---

// ## 🎯 في Go — الهيدر بيتعامل معاه من خلال `http.ResponseWriter`

// كل **Handler** في Go يأخذ:

// ```go
// func(w http.ResponseWriter, r *http.Request)
// ```

// * `w` = هو **الرد اللي هيرجعه السيرفر**.
//   من خلاله بنكتب الـ **Header** و الـ **Body** (المحتوى نفسه).
// * `r` = هو **الطلب اللي جاي من العميل**.

// ---

// ## 🧠 استخدام الـ Header في الكود

// ### 1️⃣ لإرسال نوع المحتوى (Content-Type)

// في الكود عندك مثلًا هنا:

// ```go
// w.Header().Set("Content-Type", "application/json")
// ```

// هذا السطر بيقول للمتصفح:

// > "الرد اللي هتشوفه من السيرفر هو JSON، مش HTML أو نص عادي."

// بدونه، المتصفح ممكن يعرض الـ JSON كـ **نص خام** بدون تلوين أو تفسير.

// ---

// ### 2️⃣ لإرجاع كود الحالة (Status Code)

// ```go
// w.WriteHeader(http.StatusOK)
// ```

// ده يرسل **كود الاستجابة** (زي 200 أو 404 أو 500).

// | الكود | المعنى                               |
// | ----- | ------------------------------------ |
// | `200` | OK – الطلب ناجح                      |
// | `400` | Bad Request – خطأ من العميل          |
// | `401` | Unauthorized – غير مصرح              |
// | `404` | Not Found – الصفحة غير موجودة        |
// | `500` | Internal Server Error – خطأ بالسيرفر |

// مثلاً لو حصل خطأ أثناء تكوين JSON:

// ```go
// http.Error(w, "Error generating JSON", http.StatusInternalServerError)
// ```

// ده بيرجع:

// ```http
// Status: 500 Internal Server Error
// Content-Type: text/plain
// ```

// ---

// ### 3️⃣ ممكن نضيف Headers مخصصة

// مثلاً لو عايز تضيف معلومة إضافية للـ Response:

// ```go
// w.Header().Set("X-Powered-By", "GoLang-BookingApp")
// ```

// ده يضيف Header شكله كده في الرد:

// ```http
// X-Powered-By: GoLang-BookingApp
// ```

// ---

// ## 📦 مثال عملي كامل (من كودك المعدل)

// ```go
// func (m *Respostory) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
//     response := jsonResponse{
//         OK:      true,
//         Message: "Available!",
//     }

//     // 👇 نعيّن نوع الرد (مهم جدًا)
//     w.Header().Set("Content-Type", "application/json")

//     // 👇 نعيّن كود النجاح
//     w.WriteHeader(http.StatusOK)

//     // 👇 نكتب JSON فعلاً
//     out, err := json.MarshalIndent(response, "", "    ")
//     if err != nil {
//         http.Error(w, "Error creating JSON", http.StatusInternalServerError)
//         return
//     }

//     w.Write(out)
// }
// ```

// الرد من السيرفر يكون كده:

// ```http
// HTTP/1.1 200 OK
// Content-Type: application/json
// X-Powered-By: GoLang-BookingApp

// {
//     "ok": true,
//     "message": "Available!"
// }
// ```

// ---

// ## 🧭 تلخيص سريع

// | الوظيفة                                 | الهدف                                 | مثال                 |
// | --------------------------------------- | ------------------------------------- | -------------------- |
// | `w.Header().Set("Content-Type", "...")` | يحدد نوع المحتوى (HTML, JSON, نص...)  | `"application/json"` |
// | `w.WriteHeader(statusCode)`             | يحدد كود الاستجابة (200, 404, 500...) | `http.StatusOK`      |
// | `http.Error(w, "msg", code)`            | يرسل خطأ جاهز برمز الحالة             |                      |
// | `w.Write([]byte("text"))`               | يكتب نص الرد (Body)                   | `"Hello!"`           |

// ---

// هل ترغب أن أشرح كمان **الفرق بين الهيدر في الطلب (Request Header)** والهيدر في الرد (Response Header)** مع مثال عملي باستخدام Postman أو المتصفح؟
