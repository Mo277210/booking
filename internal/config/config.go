package config

// كلمة Config (كونفيج) = اختصار لـ Configuration
// ومعناها إعدادات أو تهيئة المشروع/التطبيق.

// 🔹 في البرمجة (خصوصًا Go)

// عادةً بنعمل ملف أو باكيج اسمه config نحط فيه كل الإعدادات اللي هنحتاجها في المشروع.

// الهدف: بدل ما نكتب نفس الإعدادات في كذا مكان، نخليها في مكان واحد ونستدعيها.

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
