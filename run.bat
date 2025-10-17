@echo off
echo 🚀 Starting Booking Web Server...
cd cmd\web
REM تجاهل ملفات test عند التشغيل
go run hello_world_web_app.go routes.go middleware.go
pause
