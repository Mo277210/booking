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
