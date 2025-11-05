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

	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/handlers"
	// "github.com/bmizerany/pat"
	// 	مكتبة **bmizerany/pat** هي مكتبة Go مفتوحة المصدر تُستخدم لإنشاء **multiplexer** (mux) لطلبات HTTP بأسلوب مشابه لـ **Sinatra** في Ruby. تتيح هذه المكتبة للمطورين تعريف مسارات (routes) مع أنماط (patterns) تحتوي على متغيرات، مما يسهل التعامل مع الطلبات الديناميكية.

	// ### 📌 الاستخدام الأساسي

	// تُستخدم المكتبة بشكل أساسي لإنشاء خوادم ويب بسيطة باستخدام حزمة `net/http` في Go. تدعم المكتبة تعريف مسارات تحتوي على متغيرات، مثل `/hello/:name`، حيث يمكن الوصول إلى قيمة المتغير `name` من خلال `req.URL.Query().Get(":name")`.

	// **مثال:**

	// ```go
	// package main

	// import (
	// 	"io"
	// 	"net/http"
	// 	"github.com/bmizerany/pat"
	// 	"log"
	// )

	// func HelloServer(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, "Hello, "+req.URL.Query().Get(":name")+"!\n")
	// }

	// func main() {
	// 	m := pat.New()
	// 	m.Get("/hello/:name", http.HandlerFunc(HelloServer))

	// 	http.Handle("/", m)
	// 	err := http.ListenAndServe(":12345", nil)
	// 	if err != nil {
	// 		log.Fatal("ListenAndServe: ", err)
	// 	}
	// }
	// ```

	// في هذا المثال، يقوم الخادم بالاستماع على المنفذ 12345، وعند الوصول إلى المسار `/hello/:name`، يتم استدعاء دالة `HelloServer` التي تعرض رسالة تحتوي على اسم الشخص المرسل في الرابط.

	// ### 🔍 المزايا الرئيسية

	// * **مطابقة الأنماط (Pattern Matching):** تدعم المكتبة مطابقة الأنماط باستخدام متغيرات تبدأ بـ `:`، مما يسهل التعامل مع المسارات الديناميكية.

	// * **إعادة التوجيه التلقائي:** عند تعريف مسار ينتهي بشرطة مائلة `/`، تقوم المكتبة تلقائيًا بإعادة التوجيه إلى المسار بدون `/`، والعكس صحيح.

	// * **توافق مع `net/http`:** المكتبة متوافقة تمامًا مع حزمة `net/http` في Go، مما يسهل دمجها في التطبيقات الحالية.

	// ### 📦 التثبيت

	// لتثبيت المكتبة، يمكنك استخدام الأمر التالي:

	// ```bash
	// go get github.com/bmizerany/pat
	// ```

	// ### 📚 المراجع

	// * [المستودع الرسمي على GitHub](https://github.com/bmizerany/pat)
	// * [الوثائق على GoDoc](https://pkg.go.dev/github.com/bmizerany/pat)

	// إذا كنت بحاجة إلى مساعدة إضافية أو أمثلة تطبيقية، فلا تتردد في السؤال!
	// package main

	// import (
	// 	"fmt"
	// 	"net/http"

	// 	"github.com/bmizerany/pat"
	// )

	// func main() {
	// 	// إنشاء multiplexer باستخدام pat
	// 	mux := pat.New()

	// 	// مسار ثابت
	// 	mux.Get("/welcome", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "Welcome to our website!")
	// 	}))

	// 	// مسار ديناميكي يحتوي على متغير :name
	// 	mux.Get("/hello/:name", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		name := r.URL.Query().Get(":name")
	// 		fmt.Fprintf(w, "Hello, %s!\n", name)
	// 	}))

	// 	// مسار آخر يحتوي على متغير :id
	// 	mux.Get("/user/:id/profile", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		id := r.URL.Query().Get(":id")
	// 		fmt.Fprintf(w, "Profile page of user #%s\n", id)
	// 	}))

	// 	// تسجيل الماكينة مع net/http
	// 	http.Handle("/", mux)

	// 	fmt.Println("Server running at http://localhost:8080")
	// 	http.ListenAndServe(":8080", nil)
	// }

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	// 	تمام، هذه المكتبات تستخدم لإنشاء **HTTP router** في Go بطريقة متقدمة ومرنة أكثر من `bmizerany/pat`. دعني أوضح لك:
	// ---
	// ### 1️⃣ المكتبة الأساسية: `github.com/go-chi/chi`
	// * هي مكتبة **Router خفيفة وسريعة** لتطبيقات Go.
	// * تسمح لك بتعريف المسارات (Routes) الثابتة والديناميكية بسهولة.
	// * تدعم المجموعات (Route Groups) و **Middlewares**.
	// **مثال سريع:**
	// ```go
	// package main
	// import (
	// 	"fmt"
	// 	"net/http"
	// 	"github.com/go-chi/chi/v5"
	// )
	// func main() {
	// 	r := chi.NewRouter()
	// 	// مسار ثابت
	// 	r.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "Welcome to Chi router!")
	// 	})
	// 	// مسار ديناميكي
	// 	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	// 		name := chi.URLParam(r, "name")
	// 		fmt.Fprintf(w, "Hello, %s!\n", name)
	// 	})
	// 	http.ListenAndServe(":8080", r)
	// }
	// ```
	// ---
	// ### 2️⃣ المكتبة المساعدة: `github.com/go-chi/chi/middleware`
	// * تحتوي على **Middleware جاهزة** تسهّل التعامل مع HTTP requests.
	// * بعض الأمثلة:
	//   * `middleware.Logger` → لتسجيل جميع الطلبات في السيرفر.
	//   * `middleware.Recoverer` → منع انهيار السيرفر عند حدوث panic.
	//   * `middleware.Timeout` → لتحديد وقت انتهاء لكل طلب.
	// **مثال مع Middleware:**
	// ```go
	// package main
	// import (
	// 	"fmt"
	// 	"net/http"
	// 	"time"
	// 	"github.com/go-chi/chi/v5"
	// 	"github.com/go-chi/chi/v5/middleware"
	// )
	// func main() {
	// 	r := chi.NewRouter()
	// 	// إضافة Middleware
	// 	r.Use(middleware.Logger)      // لتسجيل كل الطلبات
	// 	r.Use(middleware.Recoverer)   // منع انهيار السيرفر عند حدوث panic
	// 	r.Use(middleware.Timeout(10 * time.Second)) // تحديد مهلة لكل طلب
	// 	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	// 		name := chi.URLParam(r, "name")
	// 		fmt.Fprintf(w, "Hello, %s!\n", name)
	// 	})
	// 	http.ListenAndServe(":8080", r)
	// }
	// ```
	// ---
	// ### 🔹 مقارنة سريعة بين `pat` و `chi`
	// | الميزة               | `pat` | `chi`            |
	// | -------------------- | ----- | ---------------- |
	// | دعم Middleware       | محدود | كامل             |
	// | الأداء               | جيد   | ممتاز ومرن       |
	// | المسارات الديناميكية | نعم   | نعم مع `{param}` |
	// | المجموعات (Groups)   | لا    | نعم              |
	// | التعقيد              | بسيط  | متوسط إلى متقدم  |
	// ---
	// إذا أحببت، أقدر أن أصنع لك **نسخة من السيرفر السابق باستخدام `chi` و Middleware مع JSON API كامل** بحيث يكون قابل للتوسع لأي مشروع حقيقي.
	// هل تريد أن أفعل ذلك؟
)

func routes(app *config.AppConfig) http.Handler {

	//Using chi for routing

	// 	mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.Home))

	mux := chi.NewRouter()

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
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", http.HandlerFunc(handlers.Repo.ReservationSummary))

	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.Availability))
	mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(handlers.Repo.AvailabilityJSON))

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))
	//Enabling static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Get("/reservations-process/{src}/{id}", handlers.Repo.AdminProcessReservation)
		mux.Get("/reservations-delete/{src}/{id}", handlers.Repo.AdminDeleteReservation)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
	 
	})
	return mux
}
