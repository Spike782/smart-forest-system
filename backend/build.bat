@echo off
echo 编译智慧林草系统...
go build -o server.exe api/server.go
if %errorlevel% equ 0 (
    echo 编译成功！
) else (
    echo 编译失败！
    pause
)
