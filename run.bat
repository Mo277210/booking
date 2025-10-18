@echo off
echo 🚀 Starting Booking Web Server...
cd cmd\web
REM تجاهل ملفات test عند التشغيل
go run ./cmd/web/hello_world_web_app.go ./cmd/web/routes.go ./cmd/web/middleware.go
pause
