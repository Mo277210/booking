package models

import "githup.com/Mo277210/booking/internal/forms"

// تمام 👌 الكود اللي كاتبه في `models/templateData.go` مظبوط جدًا.

// خليني أوضحلك:

// ---

// ### الغرض من `TemplateData`

// الـ `struct` ده معمول عشان يبقى **حاوية (Container)** للبيانات اللي هتبعتها من الـ **Handlers** للـ **Templates (HTML)**.

// مثال:

// ```go
// type TemplateData struct {
//     StringMap map[string]string   // خرايط للـ strings (مثلاً اسم مستخدم أو رسالة)
//     IntMap    map[string]int      // خرايط للأرقام
//     FloatMap  map[string]float64  // خرايط للأرقام العشرية
//     Data      map[string]interface{} // أي داتا عامة (flexible)
//     CSRFToken string              // للحماية من CSRF attacks
//     Flash     string              // رسائل نجاح
//     Warning   string              // رسائل تحذير
//     Error     string              // رسائل خطأ
// }
// ```

// ---

// ### مثال عملي:

// #### في الـ Handler:

// ```go
// func (m *Respostory) About(w http.ResponseWriter, r *http.Request) {
//     stringMap := make(map[string]string)
//     stringMap["test"] = "Hello from About handler"

//     render.RenderTemplate(w, "about", &models.TemplateData{
//         StringMap: stringMap,
//         Flash:     "This is a flash message!",
//     })
// }
// ```

// ---

// #### في الـ Template (`about.html`):

// ```html
// {{ define "content" }}
// <div class="container">
//     <h1>About Page</h1>
//     <p>This came from handler: {{ index .StringMap "test" }}</p>
//     <p>Flash Message: {{ .Flash }}</p>
// </div>
// {{ end }}
// ```

// ---

// ### النتيجة:

// لما تزور `/about` في المتصفح:

// * هيظهر نص من `StringMap["test"]`
// * وهيظهر كمان الرسالة اللي في `.Flash`

// ---

// تحب أديك مثال كمان إزاي تستفيد من **IntMap** أو **Data map** داخل التمبلت؟

//TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
 	Form     *forms.Form
}