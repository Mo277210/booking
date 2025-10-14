package render

import (
	"net/http"
	"testing"

	"githup.com/Mo277210/booking/internal/models"
)

// TestAddDefaultData يتحقق من أن دالة AddDefaultData تضيف بيانات الجلسة بشكل صحيح
func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	// الحصول على طلب يحتوي على Session
	r, err := getSession()
	if err != nil {
		t.Error("Error creating session request:", err)
	}

	// وضع قيمة في الـ Session
	session.Put(r.Context(), "flash", "123")

	// اختبار دالة AddDefaultData
	result := AddDefaultData(&td, r)

	// التحقق من أن القيمة تمت قراءتها بشكل صحيح
	if result.Flash != "123" {
		t.Error("failed: expected flash value '123' not found in session")
	}
}
func TestRenderTemplate(t *testing.T){
pathToTemplates = "../../templates"

tc,err:=CreateTemplateCache()

if err!=nil{
t.Error(err)
}
app.TemplateCache=tc
app.UseCache = true

 r , err :=getSession()
if err!=nil{
	t.Error(err)
}

var ww myWriter

	// ✅ اختبار Template موجود
RenderTemplate(&ww, r, "home.html", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser:", err)
	}

	// ✅ اختبار Template غير موجود — المفروض يُرجع خطأ
	err = RenderTemplate(&ww, r, "non-existing.html", &models.TemplateData{})
	if err == nil {
		t.Error("Expected error for non-existing template, but got none")
	}

}
// getSession تُعيد طلب جديد مع سياق الجلسة (Session context)
func getSession() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/some-url", nil)
    if err != nil {
        return nil, err
    }

    ctx := r.Context()
    ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
    r = r.WithContext(ctx)

    return r, nil
}
func TestNewlates(t *testing.T){
	NewTemplates(app)
}


func TestCreatetemplateChache(t *testing.T){
pathToTemplates = "../../templates"
_,err:=CreateTemplateCache()
if err!=nil{
	t.Error(err)
}

}