package main

//Using pat for routing
// سؤال ممتاز 👌

// ### 📌 يعني إيه **Routes** (الراوتس)؟

// كلمة **Route** معناها "مسار".
// في برمجة الويب، **Routes** = المسارات أو العناوين (URLs) اللي السيرفر بيسمع لها، وكل Route بيرتبط بـ **Handler** معين.

// ---

// ### 🔹 مثال:

// لو عندك موقع شغال على:

// ```
// http://localhost:8080
// ```

// ممكن تعرّف Routes كده:

// ```go
// mux.Get("/", HomeHandler)       // لو الزائر دخل /
// mux.Get("/about", AboutHandler) // لو الزائر دخل /about
// mux.Get("/contact", ContactHandler) // لو دخل /contact
// ```

// ---

// ### 🔧 إيه اللي بيحصل؟

// * لما المستخدم يزور `/` → يروح للـ `HomeHandler`.
// * لما يزور `/about` → يروح للـ `AboutHandler`.
// * لما يزور `/contact` → يروح للـ `ContactHandler`.

// ---

// ### 🎯 تقدر تعتبر الـ Routes زي **خريطة السيرفر**:

// * `/` = الصفحة الرئيسية.
// * `/about` = صفحة "عنّا".
// * `/login` = صفحة تسجيل الدخول.
// * `/api/users` = Endpoint بيرجع بيانات مستخدمين (لو عندك API).

// ---

// ### 🛠️ في Go مع **chi** (أو net/http):

// ```go
// r := chi.NewRouter()

// r.Get("/", HomeHandler)
// r.Get("/about", AboutHandler)

// http.ListenAndServe(":8080", r)
// ```

// هنا:

// * `r.Get` = بنسجّل Route جديد.
// * `"/"` أو `"/about"` = هو المسار (URL).
// * `HomeHandler` أو `AboutHandler` = الكود اللي يتنفذ.

// ---

// تحب أشرحلك الفرق بين **routes** و **handlers** بحيث الصورة تكمل عندك؟


import (
	"net/http"

	"githup.com/Mo277210/booking/pkg/config"
	"githup.com/Mo277210/booking/pkg/handlers"
	// "github.com/bmizerany/pat"
	"github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	
	
	//Using chi for routing

// 	mux := pat.New()

// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
// mux.Get("/about", http.HandlerFunc(handlers.Repo.Home))


mux:=	chi.NewRouter()

mux.Use(middleware.Recoverer)
//Developing our own middleware
// mux.Use(WriteToConsole)
//Creating handlers for our forms & adding CSRF Protection
mux.Use(NoSurf)
mux.Use(SessionLoad)

mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
mux.Get("/make-reservation", handlers.Repo.Reservation)
mux.Get("/majors-suite", handlers.Repo.Majors)
mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.Availability))
mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostAvailability))
mux.Get("/availability-json", http.HandlerFunc(handlers.Repo.AvailabilityJSON))

mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))
//Enabling static files
fileServer:= http.FileServer(http.Dir("./static/"))
mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

return mux
}