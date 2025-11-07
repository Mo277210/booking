@echo off
cd /d D:\visalstadio_code\heelo_ssswold\booking

REM تشغيل المشروع مباشرة بدون بناء
go run ./cmd/web/... -dbname=booking -dbuser=postgres -dbpass=1234 -dbhost=localhost -dbport=5432 -dbssl=disable -cache=false -production=false
pause
