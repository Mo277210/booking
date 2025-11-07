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
	"fmt"
	"html/template"
	"log"
	"strings"

	"strconv"
	"time"

	"net/http"

	"github.com/go-chi/chi"
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

// Repo the repository used by the handlers
var Repo *Respostory

// Respostory is the repository type
type Respostory struct {
	App *config.AppConfig
	DB  repostiory.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Respostory {
	return &Respostory{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Respostory {
	return &Respostory{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Respostory) {
	Repo = r
}

// home handler
func (m *Respostory) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home", &models.TemplateData{})
}

//about handler

func (m *Respostory) About(w http.ResponseWriter, r *http.Request) {

	//send the data to the template
	render.Template(w, r, "about", &models.TemplateData{})

}

// Reservation renders the make a reservation page and displays form
func (m *Respostory) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {

		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Respostory) PostReservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't parse reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// sd:=r.Form.Get("start_date")
	// ed:=r.Form.Get("end_date")

	// //2020-01-01 --01/02 03:04:05PM '06-0700
	// layout:="2006-01-02"
	// StartDate,err:=time.Parse(layout,sd)
	// if err!=nil {
	//     helpers.ServerError(w,err)
	//     return
	// }
	// endDate,err:=time.Parse(layout,ed)
	// if err!=nil {
	//     helpers.ServerError(w,err)
	//     return
	// }

	// roomID,err:=strconv.Atoi(r.Form.Get("room_id"))
	// if err!=nil {
	//     helpers.ServerError(w,err)
	//     return
	// }

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//reservationervation := models.Reservation{
	// 	FirstName: r.Form.Get("first_name"),
	// 	LastName:  r.Form.Get("last_name"),
	// 	Email:     r.Form.Get("email"),
	// 	Phone:     r.Form.Get("phone"),
	//     StartDate:StartDate ,
	//     EndDate:endDate,
	//     RoomID:    roomID,

	// }

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

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestrictions{

		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//send notification - first to guest
	htmlMessage := fmt.Sprintf(`
 
 <storng>Reservation Confirmation</strong><br>
    Dear %s:,<br>
    This is to confirm your reservation from %s to %s.
 `, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  template.HTML(htmlMessage),
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
	//if we reach here means the form is valid so we can put the reservation in the session ويحمل الصحفة من تاني
	// http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Respostory) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals", &models.TemplateData{})
}

// Majors renders the room page
func (m *Respostory) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Respostory) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability", &models.TemplateData{})
}

// PostAvailability handles the post request of search availability form
func (m *Respostory) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	StartDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(StartDate, endDate)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms
	// res is save in the session
	res := models.Reservation{
		StartDate: StartDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Respostory) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	//need to parse request body to get data from the form
	err := r.ParseForm()

	if err != nil {
		//can not parse the form,,so return appropriate json
		//Parsing تاريخ (Date)
		// t, err := time.Parse("2006-01-02", "2025-10-26")

		// 🔹 هنا Parse معناها:
		// حوّل النص "2025-10-26" إلى كائن وقت (time.Time) يستطيع Go التعامل معه.
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"

	StartDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(StartDate, endDate, roomID)

	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to database",
		}
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return

	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}

	//i remove the error check,since we handle all aspects of
	// the json right here

	out, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Respostory) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary
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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	m.App.Session.Remove(r.Context(), "reservation")

	render.Template(w, r, "reservation-summary", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms

func (m *Respostory) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Get(r.Context(), "reservation")

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {

		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// BookRoom takes URL parameters, assigns them to the session, and redirects to /make-reservation
func (m *Respostory) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("start")
	ed := r.URL.Query().Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// ShowLogin shows the login page
func (m *Respostory) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logging the user in
func (m *Respostory) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login", &models.TemplateData{
			Form: form,
		})
		return

	}

	email := form.Get("email")
	password := form.Get("password")

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return

	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout logs the user out
func (m *Respostory) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Respostory) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard", &models.TemplateData{})
}

// AdminAllReservations shows all reservations in admin area
func (m *Respostory) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations", &models.TemplateData{
		Data: data,
	})

}

// AdminNewReservations shows new reservations in admin area
func (m *Respostory) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows reservation details in admin area
func (m *Respostory) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	// get reservation from database

	re, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = re

	render.Template(w, r, "admin-reservations-show", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

func (m *Respostory) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateResertvation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

// AdminReservationsCalendar displays the reservation calendar
func (m *Respostory) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume that there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		// create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all the restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		
// fmt.Printf("🧩 Room %d Restrictions: %+v\n", x.ID, restrictions)

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// it's a reservation
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
			  // it's a block (owner block)
			for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
				blockMap[d.Format("2006-01-2")] = y.ID
			}
			}
		}
		
		  // 👇 أضف السطر ده هنا
    // fmt.Printf("🔹 Room ID %d | Block Map: %+v\n", x.ID, blockMap)
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calendar", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// AdminProcessReservation marks a reservation as processed
func (m *Respostory) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	err := m.DB.UpdateProcessedReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminDeleteReservation deletes a reservation
func (m *Respostory) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	err := m.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminPostReservationsCalendar handles post of reservation calendar
func (m *Respostory) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// process blocks
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {
		// Get the block map from the session. Loop through entire map, if we have an entry in the map
		// that does not exist in our posted data, and if the restriction id > 0, then it is a block we need to
		// remove.
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			// ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				// only pay attention to values > 0, and that are not in the form post
				// the rest are just placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						// delete the restriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	// now handle new blocks
	for name := range r.PostForm {
	if strings.HasPrefix(name, "add_block") {
	exploded := strings.Split(name, "_")
	if len(exploded) < 4 {
		log.Printf("⚠️ Skipping malformed block name (too short): %s -> %v", name, exploded)
		continue
	}

	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		log.Printf("⚠️ Invalid room ID in name: %s -> %v", name, exploded)
		continue
	}

	t, err := time.Parse("2006-01-2", exploded[3])
	if err != nil {
		log.Printf("⚠️ Invalid date format in name: %s -> %v", name, exploded)
		continue
	}

	err = m.DB.InsertBlockForRoom(roomID, t)
	if err != nil {
		log.Printf("❌ Error inserting block for room %d at %v: %v", roomID, t, err)
	}
}

}


	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
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

