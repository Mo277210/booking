package main
//s4---Developing our own middleware
//https://github.com/bmizerany/pat
//https://github.com/go-chi/chi?tab=readme-ov-file
//https://github.com/justinas/nosurf
//s5---Installing and setting up a sessions package
//https://github.com/alexedwards/scs
// ممتاز 👌 خليني أوضحلك ببساطة:

// 📌 Middleware يعني إيه؟

// الكلمة نفسها معناها "البرمجية الوسيطة"،
// وفي مجال الـ Web Development المقصود بيها:

// كود بيتنفذ بين الطلب (Request) اللي جاي من المتصفح والرد (Response) اللي بيرجعه السيرفر.

// 🔹 في Go مع net/http أو chi:

// الـ Middleware هو دالة (أو سلسلة دوال) بتتلف حوالين الـ Handlers.
// تقدر تستخدمه عشان:

// ✅ تسجيل (Logging) كل طلب جاي.

// ✅ التحقق من الهوية (Authentication / Authorization).

// ✅ إضافة أو تعديل بيانات على الطلب قبل ما يوصل للـ Handler.

// ✅ التحكم في الـ Response (زي إضافة Headers).
import (
	"fmt"
	"net/http"
	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from middleware")
		next.ServeHTTP(w, r)
	})
}
//NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {

csrfHandler:=nosurf.New(next)

csrfHandler.SetBaseCookie(http.Cookie{
	Name:"csrf_token",
	HttpOnly:true,
Path: "/",
Secure: app.InProduction,
SameSite: http.SameSiteLaxMode,
})
return csrfHandler
}

//SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
    return app.Session.LoadAndSave(next)
}
