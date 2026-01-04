@echo off
echo 启动智慧林草系统 Web 服务器...
cd /d %~dp0
go run api/server.go
pause

