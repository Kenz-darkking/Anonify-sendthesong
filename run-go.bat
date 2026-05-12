@echo off
setlocal
cd /d %~dp0\goapp
if not defined SENDTHESONG_DSN set SENDTHESONG_DSN=root:@tcp(127.0.0.1:3306)/sendthesong?parseTime=true
if not defined PORT set PORT=8080
"C:\Program Files\Go\bin\go.exe" run .
endlocal
